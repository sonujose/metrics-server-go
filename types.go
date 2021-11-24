package main

type httpResponse struct {
	IsSuccess bool   `json:"isSuccess"`
	Data      string `json:"data"`
	Errors    string `json:"errors"`
}

type ResponseMetricsForm struct {
	Service         string `form:"service" binding:"required"`
	Method          string `form:"method" binding:"required"`
	ResponseCode    int    `form:"responseCode" binding:"required"`
	Path            string `form:"path" binding:"required"`
	ProxyLatency    string `form:"proxyLatency" binding:"required"`
	UpstreamLatency string `form:"upstreamLatency" binding:"required"`
}
