package objects

import (
	"api_service/es"
	"api_service/heartbeat"
	"api_service/locate"
	"api_service/monit"
	"api_service/rs"
	"api_service/utils"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func post(w http.ResponseWriter, r *http.Request) {
	sTime := time.Now()
	method := r.Method
	var isok string
	defer func() {
		logger.Println("isok is", isok)
		monit.Http_request_histogram.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS, strings.Split(r.URL.Path, "/")[1], method, isok).Observe(time.Since(sTime).Seconds())
		monit.Http_request_total.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS, strings.Split(r.URL.Path, "/")[1], method, isok).Inc()
	}()
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		logger.Println(err)
		isok = "false"
		w.WriteHeader(http.StatusForbidden)
		return
	}
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		logger.Println("missing object hash in digest header")
		monit.Http_request_histogram.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS, strings.Split(r.URL.Path, "/")[1], method, "false").Observe(time.Since(sTime).Seconds())
		w.WriteHeader(http.StatusForbidden)
		isok = "false"
		return
	}
	//logger.Println("name is", name, "hash is", hash, "size is", size)
	if locate.Exist(hash) {
		err := es.AddVersion(name, hash, size)
		if err != nil {
			logger.Println(err)
			monit.Http_request_histogram.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS, strings.Split(r.URL.Path, "/")[1], method, "false").Observe(time.Since(sTime).Seconds())
			w.WriteHeader(http.StatusInternalServerError)
			isok = "false"
		} else {
			monit.Http_request_histogram.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS, strings.Split(r.URL.Path, "/")[1], method, "true").Observe(time.Since(sTime).Seconds())
			w.WriteHeader(http.StatusOK)
			isok = "true"
			//logger.Println(isok)
			w.Write([]byte("exist"))
		}
		return
	}
	//logger.Println("object is not exist")
	ds := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(ds) != rs.ALL_SHARDS {
		logger.Println("can not find enough data server")
		w.WriteHeader(http.StatusServiceUnavailable)
		isok = "false"
		return
	}
	//logger.Println(">>>>>>>>>>", name, hash, size, ds)
	stream, err := rs.NewRSResumablePutStream(ds, name, hash, size)
	if err != nil {
		logger.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		isok = "false"
		return
	}
	//logger.Println("<<<<<<<<<<", stream)
	location := fmt.Sprintf("/temp/%s", url.PathEscape(stream.ToToken()))
	//logger.Println(">>>>>>>>>>>", location)
	w.Header().Set("location", location)
	w.WriteHeader(http.StatusCreated)
	isok = "true"
	//logger.Println(isok)
	return
}
