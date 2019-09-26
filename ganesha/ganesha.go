package ganesha

import (
	"context"
	"log"
	"os"
	"os/exec"
)

const (
	nfsDaemon string = "/usr/bin/ganesha.nfsd"
)

// Ganesha manages the main nfs-ganesha process.
type Ganesha struct {
	cmd *exec.Cmd
	mgr *AdminMgr
}

// New creates a new nfs-ganesha process which can be Run and Closed.
func New(config string) *Ganesha {

	// NewAdminMgr() will error if the system DBus is not operational.  Make
	// sure DBus is running.  The AdminMgr will be used to read nfs-ganesha
	// status.
	mgr, err := NewAdminMgr()
	if err != nil {
		log.Fatal(err)
	}
	return &Ganesha{
		cmd: &exec.Cmd{
			Path: nfsDaemon,
			Args: []string{
				nfsDaemon,
				"-F",
				"-f", config,
				"-L", "/dev/stdout",
			},
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		},
		mgr: mgr,
	}
}

// Run starts the nfs-ganesha process, returning an immediate error and nil
// channel if the process cannot be started.
//
// If the process is successfully started, the returned channel will have the
// exit error pushed to it (which may be any error, but exec.ExitError is
// returned typically.)
//
// Once the process stops, the returned channel is closed.
func (g *Ganesha) Run() (<-chan error, error) {

	if err := g.cmd.Start(); err != nil {
		return nil, err
	}

	errCh := make(chan error)
	go func() {
		err := g.cmd.Wait()
		if err != nil {
			errCh <- err
			return
		}
		close(errCh)
	}()

	return errCh, nil
}

// Close sends a SIGINT to the nfs-ganesha process.
//
// Once the process has stopped, the channel returned from Run is closed.
func (g *Ganesha) Close(ctx context.Context) {
	if g.cmd != nil {
		_ = g.cmd.Process.Signal(os.Interrupt)
	}
}

// MonitorStatus listens for status updates and publishes to all status
// watchers.
func (g *Ganesha) MonitorStatus(ctx context.Context) error {
	return g.mgr.MonitorStatus(ctx)
}

// IsReady returns true if the nfs-ganesha is ready for operation, or false if a
// heartbeat was not received within the timeout period.
//
// The heartbeat message includes a boolean status field which is returned, but
// will always be set to true:
// https://github.com/nfs-ganesha/nfs-ganesha/blob/master/src/dbus/dbus_heartbeat.c#L54
//
// Heartbeats will not be sent when the server is unhealthy.
func (g *Ganesha) IsReady(ctx context.Context) bool {

	statusCh := make(chan bool)
	errCh := make(chan error)

	g.mgr.AddStatusWatcher(ctx, statusCh, errCh)
	defer g.mgr.RemoveStatusWatcher(ctx, statusCh)

	select {
	case <-ctx.Done():
		log.Printf("timed out waiting for nfs-ganesha heartbeat")
		return false
	case err := <-errCh:
		log.Printf("finished watching for nfs-ganesha heartbeats: %v", err)
		return false
	case ok := <-statusCh:
		return ok
	}
}
