package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/storageos/nfs/ganesha"
)

const (
	exportsPrefix = "storageos"
)

var exportDescriptors = map[string]IODescriptors{
	NFSv40: IODescriptors{
		Requested: prometheus.NewDesc(
			exportsPrefix+"_nfs_v40_requested_bytes_total",
			"Number of requested bytes for NFSv4.0 operations",
			[]string{"op", "name", "namespace"}, nil,
		),
		Transferred: prometheus.NewDesc(
			exportsPrefix+"_nfs_v40_transfered_bytes_total",
			"Number of transfered bytes for NFSv4.0 operations",
			[]string{"op", "name", "namespace"}, nil,
		),
		Operations: prometheus.NewDesc(
			exportsPrefix+"_nfs_v40_operations_total",
			"Number of operations for NFSv4.0",
			[]string{"op", "name", "namespace"}, nil,
		),
		Errors: prometheus.NewDesc(
			exportsPrefix+"_nfs_v40_operations_errors_total",
			"Number of operations in error for NFSv4.0",
			[]string{"op", "name", "namespace"}, nil,
		),
		Latency: prometheus.NewDesc(
			exportsPrefix+"_nfs_v40_operations_latency_seconds_total",
			"Cumulative time consumed by operations for NFSv4.0",
			[]string{"op", "name", "namespace"}, nil,
		),
		QueueWait: prometheus.NewDesc(
			exportsPrefix+"_nfs_v40_operations_queue_wait_seconds_total",
			"Cumulative time spent in rpc wait queue for NFSv4.0",
			[]string{"op", "name", "namespace"}, nil,
		),
	},
	NFSv41: IODescriptors{
		Requested: prometheus.NewDesc(
			exportsPrefix+"_nfs_v41_requested_bytes_total",
			"Number of requested bytes for NFSv4.1 operations",
			[]string{"op", "name", "namespace"}, nil,
		),
		Transferred: prometheus.NewDesc(
			exportsPrefix+"_nfs_v41_transfered_bytes_total",
			"Number of transfered bytes for NFSv4.1 operations",
			[]string{"op", "name", "namespace"}, nil,
		),
		Operations: prometheus.NewDesc(
			exportsPrefix+"_nfs_v41_operations_total",
			"Number of operations for NFSv4.1",
			[]string{"op", "name", "namespace"}, nil,
		),
		Errors: prometheus.NewDesc(
			exportsPrefix+"_nfs_v41_operations_errors_total",
			"Number of operations in error for NFSv4.1",
			[]string{"op", "name", "namespace"}, nil,
		),
		Latency: prometheus.NewDesc(
			exportsPrefix+"_nfs_v41_operations_latency_seconds_total",
			"Cumulative time consumed by operations for NFSv4.1",
			[]string{"op", "name", "namespace"}, nil,
		),
		QueueWait: prometheus.NewDesc(
			exportsPrefix+"_nfs_v41_operations_queue_wait_seconds_total",
			"Cumulative time spent in rpc wait queue for NFSv4.1",
			[]string{"op", "name", "namespace"}, nil,
		),
	},
	NFSv42: IODescriptors{
		Requested: prometheus.NewDesc(
			exportsPrefix+"_nfs_v42_requested_bytes_total",
			"Number of requested bytes for NFSv4.2 operations",
			[]string{"op", "name", "namespace"}, nil,
		),
		Transferred: prometheus.NewDesc(
			exportsPrefix+"_nfs_v42_transfered_bytes_total",
			"Number of transfered bytes for NFSv4.2 operations",
			[]string{"op", "name", "namespace"}, nil,
		),
		Operations: prometheus.NewDesc(
			exportsPrefix+"_nfs_v42_operations_total",
			"Number of operations for NFSv4.2",
			[]string{"op", "name", "namespace"}, nil,
		),
		Errors: prometheus.NewDesc(
			exportsPrefix+"_nfs_v42_operations_errors_total",
			"Number of operations in error for NFSv4.2",
			[]string{"op", "name", "namespace"}, nil,
		),
		Latency: prometheus.NewDesc(
			exportsPrefix+"_nfs_v42_operations_latency_seconds_total",
			"Cumulative time consumed by operations for NFSv4.2",
			[]string{"op", "name", "namespace"}, nil,
		),
		QueueWait: prometheus.NewDesc(
			exportsPrefix+"_nfs_v42_operations_queue_wait_seconds_total",
			"Cumulative time spent in rpc wait queue for NFSv4.2",
			[]string{"op", "name", "namespace"}, nil,
		),
	},
}

// ExportsCollector for NFS exports.
type ExportsCollector struct {
	name      string
	namespace string
	exportMgr *ganesha.ExportMgr
}

// NewExportsCollector creates a new collector for NFS exports.
//
// name and namespace should be set to the PVC name and namespace to label the
// metrics for the export.
func NewExportsCollector(name string, namespace string) ExportsCollector {
	mgr, err := ganesha.NewExportMgr()
	if err != nil {
		log.Fatal(err)
	}
	return ExportsCollector{
		name:      name,
		namespace: namespace,
		exportMgr: mgr,
	}
}

// Describe prometheus description
func (c ExportsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect IO stats for NFS exports.
//
// We expect there to be a single export per server, but GetIOStats will return
// an export record for each protocol that was used to access.  Each record will
// have the same ExportID.
//
// If we ever support multiple exports per server then the ExportID should be
// added as a label.  Since the ExportID is only used internally, we label with
// the PVC name & namespace instead so that user's can correlate with references
// they understand.
func (c ExportsCollector) Collect(ch chan<- prometheus.Metric) {

	stats, err := c.exportMgr.GetIOStats()
	if err != nil {
		log.Printf("failed to get nfs stats for exports: %v", err)
		return
	}

	for _, export := range stats.Exports {

		// Get descriptors for the export's specific NFS version.
		desc, ok := exportDescriptors[export.Name]
		if !ok {
			log.Printf("unhandled NFS version: %s", export.Name)
			continue
		}

		ch <- prometheus.MustNewConstMetric(
			desc.Requested,
			prometheus.CounterValue,
			float64(export.Read.Requested),
			"read", c.name, c.namespace)
		ch <- prometheus.MustNewConstMetric(
			desc.Transferred,
			prometheus.CounterValue,
			float64(export.Read.Transfered),
			"read", c.name, c.namespace)
		ch <- prometheus.MustNewConstMetric(
			desc.Operations,
			prometheus.CounterValue,
			float64(export.Read.Total),
			"read", c.name, c.namespace)
		ch <- prometheus.MustNewConstMetric(
			desc.Errors,
			prometheus.CounterValue,
			float64(export.Read.Errors),
			"read", c.name, c.namespace)
		ch <- prometheus.MustNewConstMetric(
			desc.Latency,
			prometheus.CounterValue,
			float64(export.Read.Latency)/1e9,
			"read", c.name, c.namespace)
		ch <- prometheus.MustNewConstMetric(
			desc.QueueWait,
			prometheus.CounterValue,
			float64(export.Read.QueueWait)/1e9,
			"read", c.name, c.namespace)

		ch <- prometheus.MustNewConstMetric(
			desc.Requested,
			prometheus.CounterValue,
			float64(export.Write.Requested),
			"write", c.name, c.namespace)
		ch <- prometheus.MustNewConstMetric(
			desc.Transferred,
			prometheus.CounterValue,
			float64(export.Write.Transfered),
			"write", c.name, c.namespace)
		ch <- prometheus.MustNewConstMetric(
			desc.Operations,
			prometheus.CounterValue,
			float64(export.Write.Total),
			"write", c.name, c.namespace)
		ch <- prometheus.MustNewConstMetric(
			desc.Errors,
			prometheus.CounterValue,
			float64(export.Write.Errors),
			"write", c.name, c.namespace)
		ch <- prometheus.MustNewConstMetric(
			desc.Latency,
			prometheus.CounterValue,
			float64(export.Write.Latency)/1e9,
			"write", c.name, c.namespace)
		ch <- prometheus.MustNewConstMetric(
			desc.QueueWait,
			prometheus.CounterValue,
			float64(export.Write.QueueWait)/1e9,
			"write", c.name, c.namespace)
	}
}
