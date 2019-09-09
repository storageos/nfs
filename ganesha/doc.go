// Package ganesha is a low-level wrapper for managing a single `nfs-ganesha`
// process within the container.
//
// It provides methods to Run, Close, and determine status with IsReady.
//
// Before starting `nfs-ganesha`, `dbus-daemon` must be running in `--system`
// mode.  You must compile `nfs-ganesha` with RPC disabled or `rpcbind` will
// also need to be running (RPC is required for NFSv3 but not NFSv4).
//
// Administrative tasks are performed by interacting with `nfs-ganesha` over
// DBus.  At the moment these actions include monitoring heartbeats to provide
// readiness status and to collect metrics.  See:
// https://github.com/nfs-ganesha/nfs-ganesha/wiki/Dbusinterface
package ganesha
