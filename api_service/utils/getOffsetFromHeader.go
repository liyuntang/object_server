package utils

import (
	"api_service/common"
	"log"
	"net/http"
	"strconv"
	"strings"
)
var logger *log.Logger
func init() {
	logger = common.WriteLog("logs/log.file")
}

// 这里解析客户端请求的head信息
func GetOffsetFromHeader(header http.Header) int64{
	byteRange := header.Get("range")
	if len(byteRange)<7 {
		return 0
	}
	if byteRange[:6] != "bytes=" {
		return 0
	}
	bytePos := strings.Split(byteRange[6:], "-")

	offset, _ := strconv.ParseInt(bytePos[0], 0, 64)
	//logger.Println("bytePos is", bytePos, "offset is", offset)
	return offset
}
