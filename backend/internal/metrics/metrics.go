package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics структура для хранения метрик Prometheus
type Metrics struct {
	HttpRequestsTotal   *prometheus.CounterVec
	HttpRequestsLatency *prometheus.HistogramVec
}

// NewMetrics создает и регистрирует метрики
func NewMetrics() *Metrics {
	return &Metrics{
		HttpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		HttpRequestsLatency: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_requests_latency_seconds",
				Help:    "HTTP request latency in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
	}
}
