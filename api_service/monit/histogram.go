package monit

import "github.com/prometheus/client_golang/prometheus"

var (
	Http_request_histogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_histogram",
			Help: "the histogram of http request",
			Buckets: []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1, 1.5, 2, 3, 4, 5},
		},
		[]string{"hostName", "endPoint", "location", "method", "isok"},
		)
)

func init()  {
	prometheus.MustRegister(Http_request_histogram)
}
