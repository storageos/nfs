package dbus

import (
	"context"
	"os"
	"os/exec"

	"github.com/godbus/dbus"
)

const (
	dbusUUIDGen string = "/usr/bin/dbus-uuidgen"
	dbusDaemon  string = "/usr/bin/dbus-daemon"
)

// DBus manages the DBus subsystem.
type DBus struct {
	cmd *exec.Cmd
}

// New creates a new DBus instance which can be Run and Closed.
func New() *DBus {
	return &DBus{
		cmd: &exec.Cmd{
			Path: dbusDaemon,
			Args: []string{
				dbusDaemon,
				"--system",
				"--nofork",
				"--nopidfile",
			},
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		},
	}
}

// Run starts the dbus-daemon process, returning an immediate error and nil
// channel if the process cannot be started.
//
// If the process is successfully started, the returned channel will have the
// exit error pushed to it (which may be any error, but exec.ExitError is
// returned typically.)
//
// Once the process stops, the returned channel is closed.
func (d *DBus) Run() (<-chan error, error) {

	if err := d.prepare(); err != nil {
		return nil, err
	}

	if err := d.cmd.Start(); err != nil {
		return nil, err
	}

	errCh := make(chan error)
	go func() {
		err := d.cmd.Wait()
		if err != nil {
			errCh <- err
			return
		}
		close(errCh)
	}()

	return errCh, nil
}

// Close sends a SIGINT to the dbus-daemon process.
//
// Once the process has stopped, the channel returned from Run is closed.
func (d *DBus) Close(ctx context.Context) {
	if d.cmd != nil {
		_ = d.cmd.Process.Signal(os.Interrupt)
	}
}

// IsReady returns true if the DBus is ready for operation.
//
// TODO: respect context.
func (d *DBus) IsReady(ctx context.Context) bool {

	conn, err := dbus.SystemBusPrivate()
	if err != nil {
		return false
	}
	defer conn.Close()

	if err := conn.Auth(nil); err != nil {
		return false
	}
	if err := conn.Hello(); err != nil {
		return false
	}

	return true
}

// prepare the system for running dbus-daemon.  Returns after command
// completion.
func (d *DBus) prepare() error {

	if err := os.MkdirAll("/run/dbus", 0755); err != nil {
		return err
	}

	// Make sure the node has a machine-id.
	idgen := &exec.Cmd{
		Path: dbusUUIDGen,
		Args: []string{
			dbusUUIDGen,
			"--ensure",
		},
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	return idgen.Run()
}
