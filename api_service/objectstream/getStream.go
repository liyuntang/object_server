package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type GetStream struct {
	reader io.Reader
}

func newGetStream(url string) (*GetStream, error) {

	//url = "http://192.168.74.98:9200/objects/aaa.1"
	//logger.Println(url)
	r, e := http.Get(url)
	if e != nil {
		logger.Println("get data is bad, err is", e)
		return nil, e
	}
	//logger.Println("get data is ok", r.StatusCode)
	//defer r.Body.Close()
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dataServer return http code %d", r.StatusCode)
	}
	//buf, _ := io.ReadAll(r.Body)
	//logger.Println("data is", string(buf))
	return &GetStream{r.Body}, nil
}

func NewGetStream(server, object string) (*GetStream, error) {
	//logger.Println(server, "===============", object)
	if server == "" || object == "" {
		return nil, fmt.Errorf("invalid server %s object %s", server, object)
	}
	url := fmt.Sprintf("http://%s/objects/%s", server, object)
	//logger.Println("url is", url)
	return newGetStream(url)
}

func (r *GetStream)Read(p []byte) (n int, err error) {
	return r.reader.Read(p)
}