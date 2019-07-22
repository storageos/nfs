// Package dbus is a simple low-level wrapper for managing a single DBus process
// within the container.
//
// It provides methods to Run, Close, and determine status with IsReady.
//
// DBus will fail to start if there is an existing unmanaged process running
// that is bound to the system bus (/var/run/dbus/system_bus_socket).
package dbus
