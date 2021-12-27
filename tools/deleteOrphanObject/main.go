package main

import (
	"deleteOrphanObject/es"
	"flag"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

var storage string
var dataServer string

func init() {
	flag.StringVar(&storage, "storage", "/Users/liyuntang/data", "storage")
	flag.StringVar(&dataServer, "dataServer", "", "dataServer address and port")
}
func main() {
	flag.Parse()
	dir := fmt.Sprintf("%s/objects/*", storage)
	files, err := filepath.Glob(dir)
	if err != nil {
		fmt.Println("scan file path of", dir, "is bad, err is", err)
		return
	}
	hashSlice := []string{}
	for _, file := range files {
		hash := strings.Split(filepath.Base(file), ".")[0]
		hashSlice = append(hashSlice, hash)
	}
	if len(hashSlice) == 0 {
		fmt.Println("there no object")
		return
	}
	// 判断hash是否在元数据表中存在，如果不存在则说明该hash可以删除了
	//hashSlice = []string{"aaa.2.fsdasafaf", "fmf1DOctL9ZOOH86HwaCVq6uh9aMxjr8Gv9c2sHKYoo="}
	sTime := time.Now()
	dHashSlice, err := es.HasHash(hashSlice)
	if err != nil {
		fmt.Println("search hash from metadata is bad, err is", err)
		return
	}
	fmt.Println(dHashSlice)
	fmt.Println("run time is", time.Since(sTime))
	// 删除object
	del(dHashSlice)
}

func del(hashSlice []string) {
	// 获取api server信息
	//fmt.Println(">>>>>>>>>>>>>>>>>>")
	for _, hash := range hashSlice {
		url := fmt.Sprintf("http://%s/objects/%s", dataServer, hash)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			fmt.Println("new request is bad, err is", err)
			return
		}
		client := http.Client{}
		resp, err := client.Do(request)
		if err != nil  || resp.StatusCode != 200{
			fmt.Println("client do is bad, err is", err)
			return
		}
		defer resp.Body.Close()
		fmt.Println("client do is ok")
	}
}


