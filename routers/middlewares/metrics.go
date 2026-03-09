package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type httpMetrics struct {
	httpRequestCount      prometheus.Counter
	httpRequestErrorCount prometheus.Counter
	httpRequestDuration   prometheus.Histogram
}

func MetricMiddleware(registry *prometheus.Registry) gin.HandlerFunc {
	metrics := &httpMetrics{
		httpRequestCount: promauto.With(registry).NewCounter(prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "The total number of http requests.",
		}),
		httpRequestErrorCount: promauto.With(registry).NewCounter(prometheus.CounterOpts{
			Name: "http_request_errors_total",
			Help: "The total number of errors in http requests.",
		}),
		httpRequestDuration: promauto.With(registry).NewHistogram(prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "The sum of duration of http requests.",
			Buckets: []float64{
				0.001,
				0.005,
				0.01,
				0.025,
				0.05,
				0.1,
				0.25,
				0.5,
				1.0,
				2.5,
				5.0,
				10.0,
			},
		}),
	}

	return func(context *gin.Context) {
		start := time.Now()
		context.Next()
		elapsed := time.Since(start)

		metrics.httpRequestCount.Inc()
		metrics.httpRequestDuration.Observe(elapsed.Seconds())
		if context.Writer.Status()%100 == 5 || context.Writer.Status()%100 == 4 {
			metrics.httpRequestErrorCount.Inc()
		}
	}
}
