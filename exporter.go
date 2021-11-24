package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RegisterMetricsExporter godoc
// register the response metrics endpoint in the gin gonic router
func registerMetricsExporter(router *gin.Engine) {

	rootPath := getEnv("ROOT_PATH", "/server")
	apiV1 := router.Group(rootPath)
	{
		apiV1.POST("/response-metrics", serverResponseMetricsExportHandler)
	}

}

// ExportserverResponseMetricsHandler godoc
// @Summary Export the requested values as metrice endpoint
// @Description
// @Tags Metrics Exporter
// @Param ResponseMetricsForm formData ResponseMetricsForm	true "Response metrics form"
// @Accept */*
// @Produce json
// @Success 200 {object} object httpResponse
// @Failure 400 {object} string httpResponse
func serverResponseMetricsExportHandler(c *gin.Context) {
	var metricsInput ResponseMetricsForm

	// Parsing request formdata to metricsresponse
	err := c.ShouldBind(&metricsInput)

	// Request validation failed due to invalid request
	// All fields in ResponseMetricsForm are mandatory
	if err != nil {
		log.Printf("Request validation failed - %v", err)
		c.AbortWithStatusJSON(400, &httpResponse{
			IsSuccess: false,
			Errors:    fmt.Sprintf("Input validation failed  - %v", err),
		})
		return
	}

	// Initiating metrics exposure
	err = exposeMetrics(metricsInput)

	// Error exposing metrics
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &httpResponse{
			IsSuccess: false,
			Errors:    fmt.Sprintf("Failed to send metrics to exporter  - %v", err),
		})
		return
	}

	// Successfully completed metrics exposure
	log.Printf("Successfully added metrics - %s-%s-%d", metricsInput.Method, metricsInput.Service, metricsInput.ResponseCode)

	c.JSON(http.StatusOK, &httpResponse{
		IsSuccess: true,
		Data:      fmt.Sprintf("Added metrics for service - %s", metricsInput.Service),
	})

}

// ExposeMetrics godoc
// @description Expose the given parameters as prometheus metrics
// @Param ResponseMetricsForm formData ResponseMetricsForm
// Creates 3 metrics
// 1) server_http_latency [histogram] {type="proxy"} 	  :  server response upstream latency
// 2) server_http_latency [histogram] {type="upstream"} :	 server response proxy latency
// 3) server_http_response [counter] 					  :  server response status counter
// Metrics Dimensions - path, method, service, statuscode
//
func exposeMetrics(metricsInput ResponseMetricsForm) error {
	log.Printf("Adding server response metrics - service:%s, path:%s, method:%s, upstreamlatency:%s, proxylatency:%s, responsecode:%d\n",
		metricsInput.Service, metricsInput.Path, metricsInput.Method, metricsInput.UpstreamLatency, metricsInput.ProxyLatency, metricsInput.ResponseCode)

	// Adding Latency for upstream and proxy with different dimensions
	serverLatencyHistogramVec.WithLabelValues(metricsInput.Path, metricsInput.Method, metricsInput.Service,
		strconv.Itoa(metricsInput.ResponseCode), "upstream").Observe(convertStringToFloat(metricsInput.UpstreamLatency))

	// Adding Latency for upstream and proxy with different dimensions
	serverLatencyHistogramVec.WithLabelValues(metricsInput.Path, metricsInput.Method, metricsInput.Service,
		strconv.Itoa(metricsInput.ResponseCode), "proxy").Observe(convertStringToFloat(metricsInput.ProxyLatency))

	// Adding response counter from server
	serverResponseCounter.WithLabelValues(metricsInput.Path, metricsInput.Method, metricsInput.Service,
		strconv.Itoa(metricsInput.ResponseCode)).Inc()

	return nil
}
