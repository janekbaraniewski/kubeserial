package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var DevicesMonitored = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "devices_monitored",
	Help: "Total number of devices monitored",
})

var DevicesMonitoredDuplicate = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "devices_monitored_duplicate",
	Help: "Total number of devices monitored DUPLICATED",
})

func init() {
	// Register custom metrics with the global prometheus registry
	metrics.Registry.MustRegister(
		DevicesMonitored,
		DevicesMonitoredDuplicate,
	)
}
