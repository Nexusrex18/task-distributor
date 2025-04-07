package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	TasksProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "worker_tasks_processed_total",
			Help: "Total number of tasks processed by worker",
		},
	)

	ProcessingDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "worker_processing_duration_seconds",
			Help:    "Duration of task processing in seconds",
			Buckets: prometheus.DefBuckets,
		},
	)
)

func Register() {
	prometheus.MustRegister(TasksProcessed)
	prometheus.MustRegister(ProcessingDuration)
}
