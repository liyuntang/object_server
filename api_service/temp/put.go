package temp

import (
	"api_service/es"
	"api_service/heartbeat"
	"api_service/locate"
	"api_service/monit"
	"api_service/rs"
	"api_service/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func put(w http.ResponseWriter, r *http.Request) {
	sTime := time.Now()
	method := r.Method
	var isok string
	defer func() {
		logger.Println("isok is", isok)
		monit.Http_request_histogram.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS, strings.Split(r.URL.Path, "/")[1], method, isok).Observe(time.Since(sTime).Seconds())
		monit.Http_request_total.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS, strings.Split(r.URL.Path, "/")[1], method, isok).Inc()
	}()

	// 解析token
	token := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 根据token创建数据流连接
	stream, err := rs.NewRSResumablePutStreamFromToken(token)
	if err != nil {
		logger.Println("new rs resumable put stream from token is bad, err is", err)
		isok = "false"
		w.WriteHeader(http.StatusForbidden)
		return
	}
	// 获取当前object已经存了多少数据，如果连接data server失败，或者object不存在则返回-1，正常情况下返回当前存储的数据量
	current := stream.CurrentSize()
	if current == -1 {
		isok = "false"
		w.WriteHeader(http.StatusNotFound)
		return
	}
	offset := utils.GetOffsetFromHeader(r.Header)
	if current != offset {
		// 数据丢了
		isok = "false"
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		return
	}
	bytes := make([]byte, rs.BLOCK_SIZE)
	// 开始存入数据
	for {
		n, e := io.ReadFull(r.Body, bytes)
		if e != nil && e != io.EOF && e != io.ErrUnexpectedEOF {
			logger.Println(e)
			isok = "false"
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		stream.Flush()
		//if current >= 10 {
		//	logger.Println("current is", current)
		//	return
		//}
		current += int64(n)
		if current > stream.Size {
			// 说明我们写入的数据比object大，这就完蛋了，数据不一致了
			stream.Commit(false)
			isok = "false"
			w.WriteHeader(http.StatusForbidden)
			return
		}
		stream.Write(bytes[:n])
		if current == stream.Size{
			// 表示数据写完了，此时可以刷数据到磁盘
			stream.Flush()
			if locate.Exist(url.PathEscape(stream.Hash)) {
				stream.Commit(false)
			} else {
				stream.Commit(true)
			}
			e = es.AddVersion(stream.Name, stream.Hash, stream.Size)
			if e != nil {
				logger.Println(e)
				isok = "false"
				w.WriteHeader(http.StatusInternalServerError)
			}
			isok = "true"
			return
		}
	}

}
