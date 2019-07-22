package ganesha

import "golang.org/x/sys/unix"

// BasicIO stores the basic statistics for NFS operations.  Each field is a
// counter that is reset when the NFS server is started or when the NFS server
// receives a request to clear the stats.
type BasicIO struct {
	Requested  uint64
	Transfered uint64
	Total      uint64
	Errors     uint64
	Latency    uint64
	QueueWait  uint64
}

// StatsBaseAnswer is the base answer to stats requests, every statistics
// related answer begins with this.
type StatsBaseAnswer struct {
	Status bool
	Error  string
	Time   unix.Timespec
}

// BasicStats is the response to IO stats call, some of the fields may not be
// filled depending of the call type and status.
type BasicStats struct {
	StatsBaseAnswer
	Read  BasicIO
	Write BasicIO
}

// ExportIOStatsList contains the stats for all exports.
type ExportIOStatsList struct {
	StatsBaseAnswer
	Exports []ExportIOStats
}

// ExportIOStats contains the stats for a single export.
//
// The Name field denotes the protocol being reported.
type ExportIOStats struct {
	ExportID uint16
	Name     string
	Read     BasicIO
	Write    BasicIO
}
