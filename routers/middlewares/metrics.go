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
	httpRequestDuration   prometheus.Counter
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
		httpRequestDuration: promauto.With(registry).NewCounter(prometheus.CounterOpts{
			Name: "http_request_duration_seconds_sum",
			Help: "The sum of duration of http requests.",
		}),
	}

	return func(context *gin.Context) {
		start := time.Now()
		context.Next()
		elapsed := time.Since(start)

		metrics.httpRequestCount.Inc()
		metrics.httpRequestDuration.Add(elapsed.Seconds())
		if context.Writer.Status()%100 == 5 {
			metrics.httpRequestErrorCount.Inc()
		}
	}
}
