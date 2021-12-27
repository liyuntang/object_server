package rs

import (
	"api_service/objectstream"
	"api_service/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func NewRSResumablePutStreamFromToken(token string) (*RSResumablePutStream, error) {

	b, e := base64.StdEncoding.DecodeString(token)
	if e != nil {
		return nil, e
	}
	var t resumableToken
	e = json.Unmarshal(b, &t)
	if e != nil {
		return nil, e
	}
	//logger.Println("uuid is", t.Uuids)
	writers := make([]io.Writer, ALL_SHARDS)
	for i := range writers {
		writers[i] = &objectstream.TempPutStream{t.Servers[i], t.Uuids[i]}
	}
	enc := NewEncoder(writers)
	return &RSResumablePutStream{&RSputStream{enc,}, &t}, nil
}
/*
CurrentSize的作用是返回object当前已经写入了多少数据
 */
func (s *RSResumablePutStream)CurrentSize() int64 {
	//se,u := getTestServerAndUuid(s.Servers, s.Uuids)
	//url := fmt.Sprintf("http://%s/temp/%s", se, u)

	//logger.Println(s.Servers)
	//logger.Println(s.Uuids)
	url := fmt.Sprintf("http://%s/temp/%s", s.Servers[0], s.Uuids[0])
	r, e := http.Head(url)
	if e != nil {
		logger.Println( e)
		return -1
	}
	defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		logger.Println(r.StatusCode)
		return -1
	}
	size := utils.GetSizeFromHeader(r.Header) * DATA_SHARDS

	if size > s.Size {
		size = s.Size
	}

	return size
}

func getTestServerAndUuid(servers, uuids []string) (string, string){
	for index, s := range servers {
		if s == "192.168.74.98:9200" {
			return s, uuids[index]
		}
	}
	return "", ""
}



