package monit

import (
	"github.com/prometheus/client_golang/prometheus"
)

var(

	// 声明控制器
	Http_request_total = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "the total number of http request",
		},
		[]string{"hostName", "endPoint", "location", "method", "isok"},
	)
)

func init() {
	// 注册
	prometheus.MustRegister(Http_request_total)
}
