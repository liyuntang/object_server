package temp

import (
	"data_server/locate"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func head(w http.ResponseWriter, r *http.Request) {
	//logger.Println(">>>>>>>>>>>>>>>>>>>>>")
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	datFile := fmt.Sprintf("%s/temp/%s.dat", locate.STORAGE_ROOT, uuid)
	//logger.Println(">>>>>>>>>>>", datFile)
	f, e := os.Open(datFile)
	if e != nil {
		logger.Println(e)
		//logger.Println("open dat file is bad")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		logger.Println(e)
		//logger.Println("get stat is bad")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//logger.Println("size is", info.Size())
	w.Header().Set("Content-Length", fmt.Sprintf("%d", info.Size()))
	w.WriteHeader(http.StatusOK)
	//w.Write([]byte(""))
	return
}
