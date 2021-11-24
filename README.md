# metrics-server-go
Helper sidecar service for exposing prometheus metrics. Application expose endpoints to update defined metrics. 

## Whats inside?
The server creates and maintains 3 metrics
 1) server_http_latency [histogram] {type="proxy"} 	  :  server response upstream latency
 2) server_http_latency [histogram] {type="upstream"} :	 server response proxy latency
 3) server_http_response [counter] 					  :  server response status counter
 Metrics Dimensions - path, method, service, statuscode

You can run the application locally using `docker-compose up` command in the application directory.