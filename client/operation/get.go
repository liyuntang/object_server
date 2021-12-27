package operation

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func Get(objectName string, version int64)  {
	apiServer := getApiServer()
	fmt.Println(">>>>>>>>>>>>>>>", apiServer)
	url := fmt.Sprintf("http://%s/objects/%s?version=%d", apiServer, objectName, version)
	//url = "http://10.10.10.116:9200/objects/NTdlZmE0MjczYjkzNjAwYWM0MjQzN2M3OGZmMzEwYzA1M2NlODU0YjQ4YzNiN2I3YjBkYTdiOTdlOGQ3ZmI4ZQ=="
	fmt.Println("url is", url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("get operater is bad, err is", err)
		return
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read body is bad, err is", err)
		return
	}

	fmt.Println(">>>>>>>>>>>>>>>>>", len(buf))
}
