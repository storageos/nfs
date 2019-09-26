package ganesha

import (
	"context"
	"sync"

	"github.com/godbus/dbus"
)

// AdminMgr is a handle to Ganesha's management interface.
type AdminMgr struct {
	conn *dbus.Conn

	// statusCh receives status updates from DBus.
	statusCh chan *dbus.Signal

	// statusWatchers holds the channels of all active status watchers.  It is
	// keyed on the channel so that it can be easily removed. It is protected by
	// mu.
	statusWatchers map[chan bool]chan error
	mu             *sync.RWMutex
}

// NewAdminMgr creates a new AdminMgr for interacting with Ganesha's management
// interface.
//
// Requires system DBus, which must be running.  Ganesha does not have to have
// registered yet.
func NewAdminMgr() (*AdminMgr, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return &AdminMgr{
		conn:           conn,
		statusCh:       make(chan *dbus.Signal),
		statusWatchers: make(map[chan bool]chan error),
		mu:             &sync.RWMutex{},
	}, nil
}

// AddStatusWatcher registers a status update subscriber channel.
func (mgr *AdminMgr) AddStatusWatcher(ctx context.Context, statusCh chan bool, errCh chan error) {
	mgr.mu.Lock()
	mgr.statusWatchers[statusCh] = errCh
	mgr.mu.Unlock()
}

// RemoveStatusWatcher deregisters a status update subscriber channel.
func (mgr *AdminMgr) RemoveStatusWatcher(ctx context.Context, ch chan bool) {
	mgr.mu.Lock()
	if _, ok := mgr.statusWatchers[ch]; ok {
		delete(mgr.statusWatchers, ch)
	}
	mgr.mu.Unlock()
}

// MonitorStatus listens for status updates and publishes to all status
// watchers.
//
// The status is received by matching the server's heartbeat signal messages
// that are published on DBus and converting to bools.
//
// Ganesha does not sent heartbeats when the server is not ready, so in practice
// only "alive/true" messages will be sent, and/or an error returned when the
// context has expired or been cancelled.
func (mgr *AdminMgr) MonitorStatus(ctx context.Context) error {

	// Match events with this signature.
	match := "type='signal',path='/org/ganesha/nfsd/heartbeat',interface='org.ganesha.nfsd.admin',member='heartbeat'"

	// Create DBus signal matcher for nfsd heartbeats.
	mgr.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, match)

	// Send status signals to statusCh
	mgr.conn.Signal(mgr.statusCh)

	for {
		select {
		case <-ctx.Done():

			// Unregister DBus signal matcher.
			mgr.conn.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, match)

			// Send error to all watchers.
			mgr.mu.RLock()
			for _, watcherErrCh := range mgr.statusWatchers {
				watcherErrCh <- ctx.Err()
			}
			mgr.mu.RUnlock()

			return ctx.Err()
		case hb := <-mgr.statusCh:

			status := hb.Body[0].(bool)

			// Send status to all watchers.
			mgr.mu.RLock()
			for watcherStatusCh := range mgr.statusWatchers {
				watcherStatusCh <- status
			}
			mgr.mu.RUnlock()
		}
	}

}
