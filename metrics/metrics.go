package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Labels used by Ganesha to identify the NFS version in use.
const (
	NFSv40 string = "NFSv40"
	NFSv41        = "NFSv41"
	NFSv42        = "NFSv42"
)

// IODescriptors contains the prometheus descriptors used to store basic
// IO stats.
type IODescriptors struct {
	Requested   *prometheus.Desc
	Transferred *prometheus.Desc
	Operations  *prometheus.Desc
	Errors      *prometheus.Desc
	Latency     *prometheus.Desc
	QueueWait   *prometheus.Desc
}

// Metrics handles metrics collection and presentation.
type Metrics struct {
	registry *prometheus.Registry
}

// New creates a new Metrics instance.
func New(name string, namespace string) *Metrics {

	reg := prometheus.NewPedanticRegistry()
	reg.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
		NewExportsCollector(name, namespace),
		NewClientsCollector(name, namespace),
	)

	return &Metrics{
		registry: reg,
	}
}

// Handler registers the http endpoint for serving metrics data.
func (s *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(s.registry, promhttp.HandlerOpts{})
}
