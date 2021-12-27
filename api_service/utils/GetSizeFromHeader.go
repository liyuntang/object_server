package utils

import (
	"net/http"
	"strconv"
)

func GetSizeFromHeader(header http.Header) int64{
	size, err := strconv.ParseInt(header.Get("Content-Length"), 0, 64)
	if err != nil {
		return 0
	}
	return size
}