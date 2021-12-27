package objects

import (
	"data_server/locate"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func delete(w http.ResponseWriter, r *http.Request) {
	// 解析hash名称
	hash := strings.Split(r.URL.EscapedPath(), "/")[2]
	//logger.Println(">>>>>>>>>>>", hash)
	// 查找以hash开头的文件名称
	objectPath := fmt.Sprintf("%s/objects/%s.*", locate.STORAGE_ROOT, hash)
	//logger.Println(">>>>>>>>>>>", objectPath)
	list, err := filepath.Glob(objectPath)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if len(list) != 1 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oAndPath := list[0]
	dir := fmt.Sprintf("%s/garbage/%s", locate.STORAGE_ROOT, time.Now().Format("2006-01-02"))
	//logger.Println("dir is", dir)
	//logger.Println("oAndPath is", oAndPath)
	objectName := filepath.Base(oAndPath)
	//logger.Println("objectName is", objectName)
	_, err = os.Stat(dir)

	if os.IsNotExist(err) {
		// 说明dir不存在，创建dir
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			logger.Println(fmt.Sprintf("mkdir %s is bad, err is %v", dir, err))
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("mkdir %s is bad, err is %v", dir, err)))
			return
		}
	}

	// 说明是其他的错误，返回错即可
	if err != nil {
		logger.Println(fmt.Sprintf("get info of %s is bad, err is %v", dir, err))
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("get info of %s is bad, err is %v", dir, err)))
		return
	}
	// 将objectmv到dir
	tagObject := fmt.Sprintf("%s/%s", dir, objectName)
	err = os.Rename(oAndPath, tagObject)
	if err != nil {
		logger.Println(fmt.Sprintf("%s rename to %s is bad, err is %v", oAndPath, tagObject, err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("%s rename to %s is bad, err is %v", oAndPath, tagObject, err)))
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

