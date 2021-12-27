package objects

import (
	"api_service/common"
	"log"
	"net/http"
	"os"
)

var logger *log.Logger
var (
	hostName,_ = os.Hostname()
)

func init() {
	logger = common.WriteLog("logs/log.file")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == http.MethodPost {
		// post返回token信息
		post(w, r)
		return
	}else if method == http.MethodPut {
		// 这里开始上传了
		//put(w, r)
		return
	} else if method == http.MethodGet {
		get(w, r)
		return
	} else if method == http.MethodDelete {
		del(w, r)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
