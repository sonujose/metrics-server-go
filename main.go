package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (

	// server Response Latency metric
	serverLatencyHistogramVec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "server_http_latency",
		Help: "Latency associated with server HTTP call",
	}, []string{"path", "method", "service", "code", "type"})

	// server http response counter metric
	serverResponseCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "server_http_response",
			Help: "Response Status for server HTTP call",
		},
		[]string{"path", "method", "service", "code"},
	)
)

func main() {

	prometheus.MustRegister(serverLatencyHistogramVec)
	prometheus.MustRegister(serverResponseCounter)

	router := gin.Default()

	// Health endpoint
	router.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "up"}) })

	// Adding metrics endpoint for to the gin-gonic server
	router.GET("/metrics", prometheusHandler())

	registerMetricsExporter(router)

	log.Printf("Starting metrics server...")

	router.Run(":" + getEnv("APP_PORT", "7006"))
}

// Initialize prometheus handler for gin-gonic
func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
