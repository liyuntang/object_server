package operation

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

func (oInfo *objectInfo)doit() string {
	//logger.Println(">>>>>>>>>>>>>>>>", oInfo)
	// 获取token
	tTime := time.Now()
	logger.Println("get token")
	err := oInfo.getToken()
	if err != nil {
		if fmt.Sprintf("%v", err) == "exist" {
			logger.Println("object is exist, run time is", time.Since(tTime))
			return "0"
		}
		logger.Println("get token bad,err is", err, "run time is", time.Since(tTime))
		return "0"
	}
	logger.Println("get token ok,run time is", time.Since(tTime), "token is", oInfo.token)
	// 获取上传进度
	url := fmt.Sprintf("http://%s%s", oInfo.apiServer, oInfo.token)
	logger.Println("get head")
	hTime := time.Now()
	contentLength, err := head(url)
	if err != nil {
		logger.Println("get head is bad, err is", err, "run time is", time.Since(hTime))
		return "0"
	}
	logger.Println("get head is ok,run time is", time.Since(hTime))
	offset, err := strconv.ParseInt(contentLength, 0, 64)
	if err != nil {
		logger.Println("parse contentLength is bad, err is", err)
		return "0"
	}
	logger.Println("已经上传了", offset, "byte数据")
	file, err := os.Open(oInfo.path)
	if err != nil {
		logger.Println("open file is bad, err is", err)
		return "0"
	}
	defer file.Close()
	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		logger.Println("seek data is bad, err is", err)
		return "0"
	}
	// 根据返回的contentLength上传数据
	logger.Println("put data to api server")
	pTime := time.Now()
	err = put(url, offset, file)
	if err != nil {
		logger.Println("put data to api server is bad, err is", err, "run time is", time.Since(pTime))
		return "0"
	}else {
		logger.Println("put data to api server is ok, run time is", time.Since(pTime), "data size is", objectSize,"MB")
		return "200"
	}
}
