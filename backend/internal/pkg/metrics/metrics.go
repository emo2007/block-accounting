package metrics

import (
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestsAccepted = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "http_request_count", // metric name
		Help: "Total requests counts",
	})

	RequestDurations = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "A histogram of the HTTP request durations in seconds.",
		Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
	})
)

func Initialize(router *chi.Mux) {
	reg := prometheus.NewRegistry()

	// Add Go module build info.
	reg.MustRegister(
		collectors.NewBuildInfoCollector(),
		collectors.NewGoCollector(),
		RequestsAccepted,
		RequestDurations,
	)

	go func() {
		router.Handle("/metrics", promhttp.HandlerFor(
			reg,
			promhttp.HandlerOpts{
				// Opt into OpenMetrics to support exemplars.
				EnableOpenMetrics: true,
			},
		))
	}()
}
