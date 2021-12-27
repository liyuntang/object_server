package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type PutSream struct {
	writer *io.PipeWriter
	c chan error
}

func NewPutStream(server, object string) *PutSream {
	reader, writer := io.Pipe()
	c := make(chan error)
	url := fmt.Sprintf("http://%s/objects/%s", server, object)
	//fmt.Println(">>>>>>>>>>>>>", url)
	go func() {
		request, _ := http.NewRequest("PUT", url, reader)
		client := http.Client{}
		r, e := client.Do(request)
		//fmt.Println("e is", e)
		if e != nil && r.StatusCode != http.StatusOK {
			e = fmt.Errorf("dataServer return http code %d", r.StatusCode)
		}
		c <- e
	}()
	return &PutSream{writer, c}
}

func (w *PutSream)Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *PutSream)Close() error {
	w.writer.Close()
	return <- w.c
}

func (w *PutSream)Commit(isok bool)  {

}