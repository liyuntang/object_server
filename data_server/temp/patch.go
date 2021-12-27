package temp

import (
	"data_server/locate"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func patch(w http.ResponseWriter, r *http.Request) {
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
	if actual > tempinfo.Size {
		os.Remove(datFile)
		os.Remove(infoFile)
		log.Println("actual size", actual, "exceeds", tempinfo.Size)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func readFromFile(uuid string) (*tempInfo, error){
	fileName := fmt.Sprintf("%s/temp/%s", locate.STORAGE_ROOT, uuid)
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf, _ := io.ReadAll(file)
	var info tempInfo
	json.Unmarshal(buf, &info)
	return &info, nil
}


