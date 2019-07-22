package ganesha

import (
	"github.com/godbus/dbus"
	"golang.org/x/sys/unix"
)

// Client Structure of the output of ShowClients dbus call.
//
// The Client field is the internal reference for the connection, typically
// `::ffff:172.17.0.1`.
//
// Whenever client traffic for a protocol is detected, the corresponding field
// for the protocol will be set to true.  When a client has multiple mounts
// using different protocols, or a single export has been remounted using a
// different protocol, mutiple protocols will be set to true for the client.
//
// LastTime is the timestamp when the response was generated.
type Client struct {
	Client   string
	NFSv3    bool
	MNTv3    bool
	NLMv4    bool
	RQUOTA   bool
	NFSv40   bool
	NFSv41   bool
	NFSv42   bool
	Plan9    bool
	LastTime unix.Timespec
}

// ClientMgr is a handle to Ganesha's ClientMgr DBus object.
//
// It's main purpose it to list clients and to retrieve per-client connection
// statistics.
type ClientMgr struct {
	dbusObject dbus.BusObject
}

// NewClientMgr returns a new ClientMgr.
func NewClientMgr() (*ClientMgr, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return &ClientMgr{
		dbusObject: conn.Object(
			"org.ganesha.nfsd",
			"/org/ganesha/nfsd/ClientMgr",
		),
	}, nil
}

// ShowClients returns Ganesha's list of client connections since the server was
// started.
func (mgr *ClientMgr) ShowClients() ([]Client, error) {

	var clients []Client
	utime := unix.Timespec{}

	if err := mgr.dbusObject.Call("org.ganesha.nfsd.clientmgr.ShowClients", 0).Store(&utime, &clients); err != nil {
		return nil, err
	}
	return clients, nil
}

// GetNFSv40IO returns basic stats for the NFSv4.0 client connection.
func (mgr *ClientMgr) GetNFSv40IO(ipaddr string) (*BasicStats, error) {
	return mgr.getBasicStats("org.ganesha.nfsd.clientstats.GetNFSv40IO", ipaddr)
}

// GetNFSv41IO returns basic stats for the NFSv4.1 client connection.
func (mgr *ClientMgr) GetNFSv41IO(ipaddr string) (*BasicStats, error) {
	return mgr.getBasicStats("org.ganesha.nfsd.clientstats.GetNFSv41IO", ipaddr)
}

func (mgr *ClientMgr) getBasicStats(method string, ipaddr string) (*BasicStats, error) {

	out := &BasicStats{}

	call := mgr.dbusObject.Call(method, 0, ipaddr)
	if call.Err != nil {
		return nil, call.Err
	}
	if !call.Body[0].(bool) {
		if err := call.Store(&out.Status, &out.Error, &out.Time); err != nil {
			return nil, err
		}
		return out, nil
	}
	if err := call.Store(&out.Status, &out.Error, &out.Time, &out.Read, &out.Write); err != nil {
		return nil, err
	}
	return out, nil
}
