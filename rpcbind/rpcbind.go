package rpcbind

import (
	"context"
	"os"
	"os/exec"
)

const (
	rpcDaemon string = "/sbin/rpcbind"
)

// RPCBind manages the main rpcbind process.
type RPCBind struct {
	cmd *exec.Cmd
}

// New creates a new rpcbind process which can be Run and Closed.
func New() *RPCBind {
	return &RPCBind{
		cmd: &exec.Cmd{
			Path: rpcDaemon,
			Args: []string{
				rpcDaemon,
				"-f",
			},
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		},
	}
}

// Run starts the rpcbind process, returning an immediate error and nil channel
// if the process cannot be started.
//
// If the process is successfully started, the returned channel will have the
// exit error pushed to it (which may be any error, but exec.ExitError is
// returned typically.)
//
// Once the process stops, the returned channel is closed.
func (g *RPCBind) Run() (<-chan error, error) {

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

// Close sends a SIGINT to the rpcbind process.
//
// Once the process has stopped, the channel returned from Run is closed.
func (g *RPCBind) Close(ctx context.Context) {
	if g.cmd != nil {
		_ = g.cmd.Process.Signal(os.Interrupt)
	}
}
