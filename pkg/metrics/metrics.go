package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	Devices = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "devices_total",
		Help: "Total number of devices",
	})
	AvailableDevices = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "devices_available",
		Help: "Total number of devices",
	})
	ReadyDevices = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "devices_ready",
		Help: "Total number of devices",
	})
)

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(
		Devices,
		AvailableDevices,
		ReadyDevices,
	)
}
