package objects

import (
	"compress/gzip"
	"crypto/sha256"
	"data_server/locate"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	file :=  getFile(strings.Split(r.URL.EscapedPath(), "/")[2])
	if file == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	sendFile(w, file)
}

func getFile(name string) string {
	//logger.Println("name is", name)
	files, _ := filepath.Glob(fmt.Sprintf("%s/objects/%s.*", locate.STORAGE_ROOT, name))
	//logger.Println("files is", files)
	if len(files) != 1 {
		return ""
	}
	//file := fmt.Sprintf("%s/objects/%s", locate.STORAGE_ROOT, name)
	//logger.Println("file is", file)
	file := files[0]
	h := sha256.New()
	sendFile(h, file)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	hash := strings.Split(filepath.Base(file), ".")[2]
	//logger.Println("d is", d)
	//logger.Println("hash is", hash)
	if d != hash {
		logger.Println("object hash mismath, remove", file)
		locate.Del(hash)
		os.Remove(file)
		return ""
	}
	//logger.Println("hash is ok")
	return file
}

func sendFile(w io.Writer, file string) {
	f, err := os.Open(file)
	if err != nil {
		logger.Println(err)
		return
	}
	defer f.Close()

	gzipStream, e := gzip.NewReader(f)
	if e != nil {
		logger.Println(e)
		return
	}

	io.Copy(w, gzipStream)
	gzipStream.Close()
}