package objectstream

import (
	"api_service/common"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type TempPutStream struct {
	Server string
	Uuid string
}
var logger *log.Logger
func init() {
	logger = common.WriteLog("logs/log.file")
}

func NewTempPutStream(server, hash string, size int64) (*TempPutStream, error) {
	url := fmt.Sprintf("http://%s/temp/%s", server, hash)
	logger.Println("url is", url, "size is", size)
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	response, err := client.Do(request)
	if response.StatusCode != http.StatusCreated {
		return nil, errors.New("http code is not 201")
	}
	defer response.Body.Close()
	uuid, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	if len(uuid) == 0{
		//logger.Println("uuid is", hash)
		return &TempPutStream{server, hash}, nil
	}
	//logger.Println("uuid is", string(uuid))
	return &TempPutStream{server, string(uuid)}, nil
}

func (w *TempPutStream)Write(p []byte) (n int, err error) {
	url := fmt.Sprintf("http://%s/temp/%s", w.Server, w.Uuid)
	request, err := http.NewRequest("PATCH", url, bytes.NewReader(p))
	if err != nil {
		return 0, err
	}
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataServer return http code %d", response.StatusCode)
	}
	return len(p), nil
}

func (w *TempPutStream) Commit(good bool) {
	method := "DELETE"
	if good {
		method = "PUT"
	}
	url := fmt.Sprintf("http://%s/temp/%s", w.Server, w.Uuid)
	request, _ := http.NewRequest(method, url, nil)
	client := http.Client{}
	response, err := client.Do(request)
	if err == nil {
		response.Body.Close()
	}
}









