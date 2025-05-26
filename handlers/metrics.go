package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	reg = prometheus.NewRegistry()

	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total HTTP requests",
		},
		[]string{"path", "method", "status"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path"},
	)

	pdfRenderDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "pdf_render_duration_seconds",
			Help:    "Duration of PDF rendering",
			Buckets: prometheus.DefBuckets,
		},
	)

	pdfRenderFailures = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "pdf_render_failures_total",
			Help: "Number of failed PDF renders",
		},
	)
)

func MetricsHandler() fiber.Handler {
	h := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
	return adaptor.HTTPHandler(h)
}

func MetricsMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start).Seconds()

		status := fmt.Sprintf("%d", c.Response().StatusCode())
		httpRequests.WithLabelValues(c.Route().Path, c.Method(), status).Inc()
		requestDuration.WithLabelValues(c.Route().Path).Observe(duration)

		return err
	}
}

func InitMetrics() {
	reg.MustRegister(httpRequests, requestDuration, pdfRenderDuration, pdfRenderFailures)
}
