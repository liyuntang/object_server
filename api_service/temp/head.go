package temp

import (
	"api_service/heartbeat"
	"api_service/monit"
	"api_service/rs"
	"fmt"
	"net/http"
	"strings"
	"time"
)

/*
客户端上传的token：1）
head函数解析出来的token：2）

1）eyJOYW1lIjoiYWFhLmZpbGUiLCJTaXplIjozMywiSGFzaCI6IjBpRXhhdHMySTVxeTRZZllBMlVkck9TM0tIaVVXU1pEYkYrUGZuT05VODQ9IiwiU2VydmVycyI6WyIxMC4xMC4xMC42Mjo5MjAwIiwiMTAuMTAuMTAuNTM6OTIwMCIsIjEwLjEwLjEwLjI0Mzo5MjAwIiwiMTAuMTAuMTAuMTg2OjkyMDAiLCIxOTIuMTY4Ljc0Ljk4OjkyMDAiLCIxMC4xMC4zMC42Mzo5MjAwIl0sIlV1aWRzIjpbIjg2NmMzNzM4LTlkZDgtNDM1ZC05ODRhLWM0ZmQ1ZjZmZDYzNyIsIjQ2NzFjZDczLWYxNDAtNGVmNy05ZmY2LTMxMzY3NDk0NDQ5OSIsIjhhNzNlNWNkLWI1ZmYtNGVmZS04MDU2LTU0MThiZmZlMjc0MSIsImU1YjgxYTdiLWY4ZjYtNDU1YS05OTdlLTAyNGRjZTM3MzBhNSIsIjVFQUZENzU2LTVEREEtNDRDMy05NTFCLThFMENCN0JENUQ1RCIsIjQ3NGEwOGM3LTg5MzQtNDZlYS04NWE5LWJhN2ExN2ZjMWZkZSJdfQ==
2）eyJOYW1lIjoiYWFhLmZpbGUiLCJTaXplIjozMywiSGFzaCI6IjBpRXhhdHMySTVxeTRZZllBMlVkck9TM0tIaVVXU1pEYkYrUGZuT05VODQ9IiwiU2VydmVycyI6WyIxMC4xMC4xMC42Mjo5MjAwIiwiMTAuMTAuMTAuNTM6OTIwMCIsIjEwLjEwLjEwLjI0Mzo5MjAwIiwiMTAuMTAuMTAuMTg2OjkyMDAiLCIxOTIuMTY4Ljc0Ljk4OjkyMDAiLCIxMC4xMC4zMC42Mzo5MjAwIl0sIlV1aWRzIjpbIjg2NmMzNzM4LTlkZDgtNDM1ZC05ODRhLWM0ZmQ1ZjZmZDYzNyIsIjQ2NzFjZDczLWYxNDAtNGVmNy05ZmY2LTMxMzY3NDk0NDQ5OSIsIjhhNzNlNWNkLWI1ZmYtNGVmZS04MDU2LTU0MThiZmZlMjc0MSIsImU1YjgxYTdiLWY4ZjYtNDU1YS05OTdlLTAyNGRjZTM3MzBhNSIsIjVFQUZENzU2LTVEREEtNDRDMy05NTFCLThFMENCN0JENUQ1RCIsIjQ3NGEwOGM3LTg5MzQtNDZlYS04NWE5LWJhN2ExN2ZjMWZkZSJdfQ==
 */
func head(w http.ResponseWriter, r *http.Request) {
	sTime := time.Now()
	method := r.Method
	var isok string
	defer func() {
		logger.Println("isok is", isok)
		monit.Http_request_histogram.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS, strings.Split(r.URL.Path, "/")[1], method, isok).Observe(time.Since(sTime).Seconds())
		monit.Http_request_total.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS, strings.Split(r.URL.Path, "/")[1], method, isok).Inc()
	}()
	//logger.Println("=======", r.URL.EscapedPath())
	// 这里解析出来的token是编码后的值
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	//logger.Println(">>>>>>>", token)
	// 这里的token跟客户端传过来的token是相同的
	stream, e := rs.NewRSResumablePutStreamFromToken(token)
	if e != nil {
		logger.Println(e)
		// 403
		w.WriteHeader(http.StatusForbidden)
		isok = "false"
		return
	}
	current := stream.CurrentSize()
	//logger.Println("current is", current)
	if current == -1 {
		// 404
		isok = "false"
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("content-length", fmt.Sprintf("%d", current))
	isok = "true"
	return
}
