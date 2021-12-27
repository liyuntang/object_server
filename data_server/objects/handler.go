package objects

import (
	"data_server/common"
	"log"
	"net/http"
)
var logger *log.Logger

func init() {
	logger = common.WriteLog("logs/log.file")
}
func Handler(w http.ResponseWriter, r *http.Request) {
	method := r.Method

	if method == http.MethodGet {
		//logger.Println("get..............")
		get(w, r)
	} else if method == http.MethodDelete {
		//logger.Println(">>>>>>>>>>>>>>>>>>> delete")
		delete(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
