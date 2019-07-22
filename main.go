package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/storageos/nfs/dbus"
	"github.com/storageos/nfs/ganesha"
	"github.com/storageos/nfs/health"
	"github.com/storageos/nfs/http"
	"github.com/storageos/nfs/metrics"
	"github.com/storageos/nfs/rpcbind"
)

const (
	name            string = "StorageOS NFS"
	healthEndpoint         = "/healthz"
	metricsEndpoint        = "/metrics"
)

const (
	ganeshaConfigEnvVar  string = "GANESHA_CONFIGFILE"
	listenAddrEnvVar     string = "LISTEN_ADDR"
	nameEnvVar           string = "NAME"
	namespaceEnvVar      string = "NAMESPACE"
	disableMetricsEnvVar string = "DISABLE_METRICS"
)

func main() {

	ganeshaConfig := os.Getenv(ganeshaConfigEnvVar)
	if ganeshaConfig == "" {
		log.Fatalf("ganesha config file must be specified with %s env var", ganeshaConfigEnvVar)
	}

	listenAddr := os.Getenv(listenAddrEnvVar)
	if listenAddr == "" {
		listenAddr = ":80"
	}

	dm := os.Getenv(disableMetricsEnvVar)
	if dm == "" {
		dm = "false"
	}
	disableMetrics, err := strconv.ParseBool(dm)
	if err != nil {
		log.Fatalf("%s env var value must be true or false/empty/unset", disableMetricsEnvVar)
	}

	// All processes should start and be ready within the context timeout.  Can
	// be extended as needed, but 10 seconds should be plenty.
	startCtx, startCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer startCancel()

	// Start rpcbind
	rpc := rpcbind.New()
	rpcErrCh, err := rpc.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Start DBus.
	bus := dbus.New()
	dbusErrCh, err := bus.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Wait for Dbus to be operational.
	if err := waitForReady(startCtx, bus.IsReady); err != nil {
		log.Fatal(err)
	}

	// Start Ganesaha.
	nfs := ganesha.New(ganeshaConfig)
	nfsErrCh, err := nfs.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Wait for Ganesha to report it is ready.
	if err := waitForReady(startCtx, nfs.IsReady); err != nil {
		log.Fatal(err)
	}

	// Start HTTP server.
	srv := http.New(listenAddr, name)
	srv.RegisterHandler("Index", "/", srv.Handler())

	httpErrCh, err := srv.Run()
	if err != nil {
		log.Fatal(err)
	}

	// Register health endpoint.
	srv.RegisterHandler("Health", healthEndpoint, health.New(nfs).Handler())

	// Register metrics endpoints if not explicitly disabled.
	if !disableMetrics {
		log.Printf("enabling prometheus endpoint on http://%s/metrics", listenAddr)
		metrics := metrics.New(os.Getenv(nameEnvVar), os.Getenv(namespaceEnvVar))
		srv.RegisterHandler("Metrics", metricsEndpoint, metrics.Handler())
	}

	var stopCh = make(chan os.Signal)
	signal.Notify(stopCh, os.Interrupt)

	select {
	case <-stopCh:
		log.Print("shutdown requested")
	case err := <-rpcErrCh:
		log.Printf("rpcbind daemon stopped: %v", err)
	case err := <-dbusErrCh:
		log.Printf("dbus daemon stopped: %v", err)
	case err := <-nfsErrCh:
		log.Printf("nfs server stopped: %v", err)
	case err := <-httpErrCh:
		log.Printf("http server stopped: %v", err)
	}

	// Graceful stop
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer stopCancel()

	nfs.Close(stopCtx)
	srv.Close(stopCtx)
	rpc.Close(stopCtx)
	bus.Close(stopCtx)

	log.Print("graceful shutdown completed")

}

// waitForReady waits for readyFunc to return true or the context to expire.
//
// Calls to readyFunc() are intended to be inexpensive, hence the minimal delay
// and no backoff, optimising for bringing services online as quickly as
// possible.
//
// Where a readyFunc() is expensive, it should introduce its own delay/backoff
// to reduce load.
func waitForReady(ctx context.Context, readyFunc func(ctx context.Context) bool) error {

	timer := time.NewTicker(100 * time.Millisecond)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			if ok := readyFunc(ctx); ok {
				return nil
			}
		}
	}
}
