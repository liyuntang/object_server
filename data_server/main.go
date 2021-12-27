package main

import (
	"data_server/heartbeat"
	"data_server/locate"
	"data_server/objects"
	"data_server/temp"
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	flag.Parse()
	locate.CollectObjects()
	go heartbeat.StartHeartbeat()
	go locate.StartLocate()
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/temp/", temp.Handler)
	log.Fatal(http.ListenAndServe(heartbeat.LISTEN_ADDRESS, nil))
}
