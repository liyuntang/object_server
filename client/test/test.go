package test

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

func Test() {
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/check", monit(check))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func check(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3*time.Second)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
	return
}
var webRequestTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "web_request_total",
		Help: "number of hello requests in toal",
	},
	[]string{"method", "endpoint"},
	)

var webRequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "web_request_duration",
		Help: "web request duration",
		Buckets: []float64{0.1, 0.3, 0.5, 0.7, 0.9, 1},
	},
	[]string{"method", "endpoint"},
	)

func init() {
	prometheus.Register(webRequestTotal)
	prometheus.Register(webRequestDuration)
}

func h(w http.ResponseWriter, r *http.Request) {

}
func monit(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		sTime := time.Now()
		duration := time.Since(sTime)
		h(writer, request)
		webRequestTotal.With(prometheus.Labels{"method":request.Method, "endpint":request.URL.Path}).Inc()
		webRequestDuration.With(prometheus.Labels{"method":request.Method, "endpoint":request.URL.Path}).Observe(duration.Seconds())
	}
}