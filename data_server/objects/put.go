package objects

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(">>>>>>>>>>>>>>>>>>>>>>>> put")
	object := fmt.Sprintf("%s/objects/%s", os.Getenv("STORAGE_ROOT"), strings.Split(r.URL.EscapedPath(), "/")[2])
	file, err := os.Create(object)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	io.Copy(file, r.Body)
	return
}
