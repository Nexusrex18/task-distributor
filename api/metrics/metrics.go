package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	TasksSubmitted = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "api_tasks_submitted_total",
			Help: "Total number of tasks submitted via API",
		},
	)
	RequestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "api_request_duration_seconds",
			Help:    "Duration of API requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func Register() {
	prometheus.MustRegister(TasksSubmitted)
	prometheus.MustRegister(RequestDuration)
}