package monit

import "github.com/prometheus/client_golang/prometheus"

var (
	Object_total = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "object_total",
			Help: "the number of objects",
		},
		[]string{"hostName", "endPoint"},
		)
)

func init() {
	prometheus.MustRegister(Object_total)
}

