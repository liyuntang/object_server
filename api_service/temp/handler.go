package temp

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
	m := r.Method
	if m == http.MethodHead {
		head(w, r)
		return
	} else if m == http.MethodPut {
		put(w, r)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

























