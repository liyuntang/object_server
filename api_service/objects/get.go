package objects

import (
	"api_service/es"
	"api_service/heartbeat"
	"api_service/locate"
	"api_service/monit"
	"api_service/rs"
	"api_service/utils"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func get(w http.ResponseWriter, r *http.Request) {
	sTime := time.Now()
	method := r.Method
	var isok2 string
	defer func() {
		logger.Println("isok is", isok2)
		monit.Http_request_histogram.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS, strings.Split(r.URL.Path, "/")[1], method, isok2).Observe(time.Since(sTime).Seconds())
		monit.Http_request_total.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS, strings.Split(r.URL.Path, "/")[1], method, isok2).Inc()
	}()
	//logger.Println(">>>>>>>>>>>>>>>>> get ")
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 获取version信息
	versionId := r.URL.Query()["version"]
	//logger.Println("version id is", versionId)
	var version int64
	var e error
	var isok bool
	if len(versionId) != 0 {
		version, e = strconv.ParseInt(versionId[0], 0, 64)
		if e != nil {
			isok2 = "false"
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		// 如果url中没有version参数说明获取该对象的最新版本，如果最新版本的size=0或者hash=""则说明该对象已经被删除，返回404即可
		isok, version, e = es.CheckObjectISExsit(name)
		//logger.Println(isok, version, e)
		if !isok {
			// 说明对象不存在
			isok2 = "false"
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	//logger.Println(name, versionId, version)
	// 根据name和version从元数据服务中获取对象信息
	meta, err := es.GetMetaData(name, version)
	if err != nil {
		//logger.Println("search meta data is bad, err is", err)
		isok2 = "false"
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if meta.Hash == "" {
		isok2 = "false"
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//logger.Println("meta is", meta)
	stream, e := GetStream(meta.Hash, meta.Size)
	if e != nil {
		log.Println(e)
		isok2 = "false"
		w.WriteHeader(http.StatusNotFound)
		return
	}
	offset := utils.GetOffsetFromHeader(r.Header)
	if offset != 0 {
		stream.Seek(offset, io.SeekCurrent)
		w.Header().Set("content-range", fmt.Sprintf("bytes=%d-%d/%d", offset, meta.Size-1, meta.Size))
		w.WriteHeader(http.StatusPartialContent)
	}
	acceptGzip := false
	encoding := r.Header["Accept-Encoding"]
	for i := range encoding {
		if encoding[i] == "gzip" {
			acceptGzip = true
			break
		}
	}
	if acceptGzip {
		w.Header().Set("content-encoding", "gzip")
		w2 := gzip.NewWriter(w)
		io.Copy(w2, stream)
		w2.Close()
	} else {
		io.Copy(w, stream)
	}
	stream.Close()
	isok2 = "true"
}

func GetStream(hash string, size int64) (*rs.RsGetStream, error) {
	locateInfo := locate.Locate(hash)
	//logger.Println("hash is", hash, "locateInfo is", locateInfo)
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("objects %s locate fail, result %v", hash, locateInfo)
	}
	dataServers := make([]string, 0)
	if len(locateInfo) != rs.ALL_SHARDS {
		dataServers = heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS-len(locateInfo), locateInfo)
	}
	return rs.NewRsGetStream(locateInfo, dataServers, hash, size)
}
