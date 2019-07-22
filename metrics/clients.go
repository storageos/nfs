package metrics

import (
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/storageos/nfs/ganesha"
)

var clientsPrefix = "storageos_clients"

var (
	clientsNfsV40RequestedDesc = prometheus.NewDesc(
		clientsPrefix+"_nfs_v40_requested_bytes_total",
		"Number of requested bytes for NFSv4.0 operations",
		[]string{"op", "name", "namespace", "clientip"}, nil,
	)
	clientsNfsV40TransferedDesc = prometheus.NewDesc(
		clientsPrefix+"_nfs_v40_transfered_bytes_total",
		"Number of transfered bytes for NFSv4.0 operations",
		[]string{"op", "name", "namespace", "clientip"}, nil,
	)
	clientsNfsV40OperationsDesc = prometheus.NewDesc(
		clientsPrefix+"_nfs_v40_operations_total",
		"Number of operations for NFSv4.0",
		[]string{"op", "name", "namespace", "clientip"}, nil,
	)
	clientsNfsV40ErrorsDesc = prometheus.NewDesc(
		clientsPrefix+"_nfs_v40_operations_errors_total",
		"Number of operations in error for NFSv4.0",
		[]string{"op", "name", "namespace", "clientip"}, nil,
	)
	clientsNfsV40LatencyDesc = prometheus.NewDesc(
		clientsPrefix+"_nfs_v40_operations_latency_seconds_total",
		"Cumulative time consumed by operations for NFSv4.0",
		[]string{"op", "name", "namespace", "clientip"}, nil,
	)
	clientsNfsV40QueueWaitDesc = prometheus.NewDesc(
		clientsPrefix+"_nfs_v40_operations_queue_wait_seconds_total",
		"Cumulative time spent in rpc wait queue for NFSv4.0",
		[]string{"op", "name", "namespace", "clientip"}, nil,
	)
	clientsNfsV41RequestedDesc = prometheus.NewDesc(
		clientsPrefix+"_nfs_v41_requested_bytes_total",
		"Number of requested bytes for NFSv4.1 operations",
		[]string{"op", "name", "namespace", "clientip"}, nil,
	)
	clientsNfsV41TransferedDesc = prometheus.NewDesc(
		clientsPrefix+"_nfs_v41_transfered_bytes_total",
		"Number of transfered bytes for NFSv4.1 operations",
		[]string{"op", "name", "namespace", "clientip"}, nil,
	)
	clientsNfsV41OperationsDesc = prometheus.NewDesc(
		clientsPrefix+"_nfs_v41_operations_total",
		"Number of operations for NFSv4.1",
		[]string{"op", "name", "namespace", "clientip"}, nil,
	)
	clientsNfsV41ErrorsDesc = prometheus.NewDesc(
		clientsPrefix+"_nfs_v41_operations_errors_total",
		"Number of operations in error for NFSv4.1",
		[]string{"op", "name", "namespace", "clientip"}, nil,
	)
	clientsNfsV41LatencyDesc = prometheus.NewDesc(
		clientsPrefix+"_nfs_v41_operations_latency_seconds_total",
		"Cumulative time consumed by operations for NFSv4.1",
		[]string{"op", "name", "namespace", "clientip"}, nil,
	)
	clientsNfsV41QueueWaitDesc = prometheus.NewDesc(
		clientsPrefix+"_nfs_v41_operations_queue_wait_seconds_total",
		"Cumulative time spent in rpc wait queue for NFSv4.1",
		[]string{"op", "name", "namespace", "clientip"}, nil,
	)
)

// ClientsCollector Collector for ganesha clients.
type ClientsCollector struct {
	name      string
	namespace string
	clientMgr *ganesha.ClientMgr
}

// NewClientsCollector creates a new collector.
func NewClientsCollector(name string, namespace string) ClientsCollector {

	mgr, err := ganesha.NewClientMgr()
	if err != nil {
		log.Fatal(err)
	}
	return ClientsCollector{
		name:      name,
		namespace: namespace,
		clientMgr: mgr,
	}
}

// Describe prometheus description
func (c ClientsCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect do the actual job
func (c ClientsCollector) Collect(ch chan<- prometheus.Metric) {

	clients, err := c.clientMgr.ShowClients()
	if err != nil {
		log.Printf("failed to get nfs client list: %v", err)
		return
	}

	for _, client := range clients {

		// Only NFSv40 and NFSv41 client metrics are supported by Ganesha.

		if client.NFSv40 {

			stats, err := c.clientMgr.GetNFSv40IO(client.Client)
			if err != nil {
				log.Printf("failed to get nfs 4.0 stats for client: %v", err)
				continue
			}

			ch <- prometheus.MustNewConstMetric(
				clientsNfsV40RequestedDesc,
				prometheus.CounterValue,
				float64(stats.Read.Requested),
				"read", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV40TransferedDesc,
				prometheus.CounterValue,
				float64(stats.Read.Transfered),
				"read", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV40OperationsDesc,
				prometheus.CounterValue,
				float64(stats.Read.Total),
				"read", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV40ErrorsDesc,
				prometheus.CounterValue,
				float64(stats.Read.Errors),
				"read", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV40LatencyDesc,
				prometheus.CounterValue,
				float64(stats.Read.Latency)/1e9,
				"read", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV40QueueWaitDesc,
				prometheus.CounterValue,
				float64(stats.Read.QueueWait)/1e9,
				"read", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV40RequestedDesc,
				prometheus.CounterValue,
				float64(stats.Write.Requested),
				"write", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV40TransferedDesc,
				prometheus.CounterValue,
				float64(stats.Write.Transfered),
				"write", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV40OperationsDesc,
				prometheus.CounterValue,
				float64(stats.Write.Total),
				"write", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV40ErrorsDesc,
				prometheus.CounterValue,
				float64(stats.Write.Errors),
				"write", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV40LatencyDesc,
				prometheus.CounterValue,
				float64(stats.Write.Latency)/1e9,
				"write", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV40QueueWaitDesc,
				prometheus.CounterValue,
				float64(stats.Write.QueueWait)/1e9,
				"write", c.name, c.namespace, client.Client)
		}
		if client.NFSv41 {

			stats, err := c.clientMgr.GetNFSv41IO(client.Client)
			if err != nil {
				log.Printf("failed to get nfs 4.1 stats for client: %v", err)
				continue
			}

			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41RequestedDesc,
				prometheus.CounterValue,
				float64(stats.Read.Requested),
				"read", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41TransferedDesc,
				prometheus.CounterValue,
				float64(stats.Read.Transfered),
				"read", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41OperationsDesc,
				prometheus.CounterValue,
				float64(stats.Read.Total),
				"read", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41ErrorsDesc,
				prometheus.CounterValue,
				float64(stats.Read.Errors),
				"read", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41LatencyDesc,
				prometheus.CounterValue,
				float64(stats.Read.Latency)/1e9,
				"read", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41QueueWaitDesc,
				prometheus.CounterValue,
				float64(stats.Read.QueueWait)/1e9,
				"read", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41RequestedDesc,
				prometheus.CounterValue,
				float64(stats.Write.Requested),
				"write", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41TransferedDesc,
				prometheus.CounterValue,
				float64(stats.Write.Transfered),
				"write", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41OperationsDesc,
				prometheus.CounterValue,
				float64(stats.Write.Total),
				"write", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41ErrorsDesc,
				prometheus.CounterValue,
				float64(stats.Write.Errors),
				"write", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41LatencyDesc,
				prometheus.CounterValue,
				float64(stats.Write.Latency)/1e9,
				"write", c.name, c.namespace, client.Client)
			ch <- prometheus.MustNewConstMetric(
				clientsNfsV41QueueWaitDesc,
				prometheus.CounterValue,
				float64(stats.Write.QueueWait)/1e9,
				"write", c.name, c.namespace, client.Client)
		}
	}
}
