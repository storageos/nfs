package ganesha

import (
	"github.com/godbus/dbus"
)

// ExportMgr is a handle to Ganesha's DBus ExportMgr object.
//
// It can be used to retrieve per-export protocol statistics.
type ExportMgr struct {
	dbusObject dbus.BusObject
}

// NewExportMgr Get a new ExportMgr
func NewExportMgr() (*ExportMgr, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	return &ExportMgr{
		dbusObject: conn.Object(
			"org.ganesha.nfsd",
			"/org/ganesha/nfsd/ExportMgr",
		),
	}, nil
}

// GetIOStats returns the basic IO stats for all exports.
func (mgr *ExportMgr) GetIOStats() (*ExportIOStatsList, error) {

	out := &ExportIOStatsList{}

	call := mgr.dbusObject.Call("org.ganesha.nfsd.exportstats.GetNFSIO", 0)
	if call.Err != nil {
		return nil, call.Err
	}

	if !call.Body[0].(bool) {
		if err := call.Store(&out.Status, &out.Error, &out.Time); err != nil {
			return nil, err
		}
		return out, nil
	}

	if err := call.Store(&out.Status, &out.Error, &out.Time, &out.Exports); err != nil {
		return nil, err
	}
	return out, nil
}
