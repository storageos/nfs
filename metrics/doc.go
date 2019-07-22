// Package metrics provides an interface between NFS Ganesha's statistics and
// Prometheus, allowing metrics to be scaped over an HTTP endpoint on request.
//
// Metrics for exports and client connections are reported on a per-protocol
// basis.
package metrics
