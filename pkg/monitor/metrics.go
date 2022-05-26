package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	DevicesAvailable = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "devices_available",
		Help: "Total number of devices available",
	})
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(
		DevicesAvailable,
	)
}
