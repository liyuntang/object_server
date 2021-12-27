package temp

import (
	"data_server/locate"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	datFile := fmt.Sprintf("%s/temp/%s.dat", locate.STORAGE_ROOT, uuid)
	f, e := os.Open(datFile)
	if e != nil {
		logger.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer f.Close()
	io.Copy(w, f)
}
