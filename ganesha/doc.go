// Package ganesha is a low-level wrapper for managing a single `nfs-ganesha`
// process within the container.
//
// It provides methods to Run, Close, and determine status with IsReady.
//
// Before starting `nfs-ganesha`, `dbus-daemon` must be running in `--system`
// mode, and `rpcbind` must also be started.  `nfs-ganesha` seems to require
// that `rpcbind` is running even though it is not used for NFSv4:
// https://github.com/nfs-ganesha/nfs-ganesha/issues/114
//
// Administrative tasks are performed by interacting with `nfs-ganesha` over
// DBus.  At the moment these actions include monitoring heartbeats to provide
// readiness status and to collect metrics.  See:
// https://github.com/nfs-ganesha/nfs-ganesha/wiki/Dbusinterface
package ganesha
