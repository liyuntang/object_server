package objects

import (
	"api_service/es"
	"api_service/heartbeat"
	"api_service/monit"
	"net/http"
	"strings"
	"time"
)
/*
del的作用是删除对象，流程如下：
	1、调用es.PutMetadata插入一条新数据，version+=1，size=0，hash=""
 */
func del(w http.ResponseWriter, r *http.Request) {
	sTime := time.Now()
	method := r.Method
	var isok string
	defer func() {
		logger.Println("isok is", isok)
		monit.Http_request_histogram.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS, strings.Split(r.URL.Path, "/")[1], method, isok).Observe(time.Since(sTime).Seconds())
		monit.Http_request_total.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS, strings.Split(r.URL.Path, "/")[1], method, isok).Inc()
	}()

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 获取该对象最新版本
	version, err := es.SearchLatestVersion(name)
	if err != nil {
		isok = "false"
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// 计算version+=1然后写入
	err = es.PutMetadata(name, version+1, 0, "")
	if err != nil {
		isok = "false"
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	isok = "true"
	return
}














