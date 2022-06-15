package webhooks

import (
	"github.com/prometheus/client_golang/prometheus"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
)

var (
	PodsHandled = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "pods_handled",
		Help: "Number of pods handled by webhook",
	})

	InjectedCommands = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "injected_commands_total",
		Help: "Number of injected commands",
	})
)

func init() {
	metrics.Registry.MustRegister(
		PodsHandled,
	)
}
