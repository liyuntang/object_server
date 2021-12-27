package main

import (
	"api_service/heartbeat"
	"api_service/locate"
	"api_service/objects"
	"api_service/temp"
	"flag"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func main() {
	flag.Parse()
	go func() {
		for {
			heartbeat.ChooseRandomDataServers(3, nil)
			time.Sleep(time.Second)
		}

	}()

	go heartbeat.ListenHeartbeat()

	// 向online通道注册自己
	go heartbeat.Regist()
	//

	// prometheus
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/temp/", temp.Handler)
	http.HandleFunc("/objects/", objects.Handler)
	// locate
	http.HandleFunc("/locate/", locate.Handler)
	//http.HandleFunc("/versions/", version.Handler)
	log.Fatal(http.ListenAndServe(heartbeat.LISTEN_ADDRESS, nil))
}
