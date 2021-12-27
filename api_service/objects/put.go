package objects

import (
	"api_service/es"
	"api_service/heartbeat"
	"api_service/locate"
	"api_service/rs"
	"api_service/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	// 获取对象的hash值,这里的hash值是客户端传过来的
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		logger.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//logger.Println(">>>>>>>>>>>>", hash)
	//获取对象大小
	size := utils.GetSizeFromHeader(r.Header)

	// 用hash值代替对象名字
	c, e := storeObject(r.Body, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		return
	}
	if c != http.StatusOK {
		w.WriteHeader(c)
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	// 元数据入库
	//
	err := es.AddVersion(name, hash, size)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

func storeObject(r io.Reader, hash string, size int64) (int, error){
	// 判断对象是否存在，如果存在则直接返回ok，记录元数据即可
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}
	//logger.Println("hash is", hash)
	//logger.Println("urlh is", url.PathEscape(hash))
	// putStream传入的是转译后的hash
	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}
	reader := io.TeeReader(r, stream)
	d := utils.CalculateHash(reader)
	if d != hash {
		stream.Commit(false)
		logger.Println("object hash is mismatch")
		return http.StatusBadRequest, fmt.Errorf("object hash is mismatch")
	}
	stream.Commit(true)
	return http.StatusOK, nil
}

func putStream(hash string, size int64) (*rs.RSputStream, error) {
	servers := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	if len(servers) != rs.ALL_SHARDS {
		logger.Println("can not find enough dataServer", servers)
		return nil, fmt.Errorf("can not find enough dataServer")
	}
	//logger.Println(">>>", hash)
	return rs.NewRSputStream(servers, hash, size)
}