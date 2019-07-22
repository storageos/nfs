package ganesha

import (
	"context"

	"github.com/godbus/dbus"
)

// AdminMgr is a handle to Ganesha's management interface.
type AdminMgr struct {
	conn *dbus.Conn
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
		conn: conn,
	}, nil
}

// StreamStatus streams the NFS server's status to the out channel.
//
// The status is received by watching for the server's heartbeat messages that
// are published on DBus and converting to bools.
//
// Ganesha does not sent heartbeats when the server is not ready, so in practice
// only "alive/true" messages will be sent, and/or an error returned when the
// context has expired or been cancelled.
func (mgr *AdminMgr) StreamStatus(ctx context.Context, out chan bool) error {

	// Create DBus watcher for nfsd heartbeats.
	match := "type='signal',path='/org/ganesha/nfsd/heartbeat',interface='org.ganesha.nfsd.admin',member='heartbeat'"
	mgr.conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, match)

	in := make(chan *dbus.Signal)
	mgr.conn.Signal(in)

	for {
		select {
		case <-ctx.Done():
			mgr.conn.BusObject().Call("org.freedesktop.DBus.RemoveMatch", 0, match)
			return ctx.Err()
		case hb := <-in:
			out <- hb.Body[0].(bool)
		}
	}

}
