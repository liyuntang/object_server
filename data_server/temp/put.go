package temp

import (
	"compress/gzip"
	"data_server/locate"
	"data_server/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, err := readFromFile(uuid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	infoFile := fmt.Sprintf("%s/temp/%s", locate.STORAGE_ROOT, uuid)
	datFile := fmt.Sprintf("%s.dat", infoFile)
	f, e := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND, 0644)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, err = io.Copy(f, r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	info, err := f.Stat()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	actual := info.Size()
	os.Remove(infoFile)
	if actual != tempinfo.Size {
		os.Remove(datFile)
		log.Println("actual size", actual, "exceeds", tempinfo.Size)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	commitTempObject(datFile, tempinfo)
}

func commitTempObject(datFile string, tempinfo *tempInfo) {
	f, _ := os.Open(datFile)
	defer f.Close()
	d := url.PathEscape(utils.CalculateHash(f))
	f.Seek(0, io.SeekStart)

	targFile := fmt.Sprintf("%s/objects/%s.%s", locate.STORAGE_ROOT, tempinfo.Name, d)
	w, _ := os.Create(targFile)
	w2 := gzip.NewWriter(w)
	io.Copy(w2, f)
	w2.Close()

	os.Remove(datFile)
	locate.Add(tempinfo.hash(), tempinfo.id())
}
