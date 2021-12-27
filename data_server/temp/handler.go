package temp

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
	m := r.Method
	if m == http.MethodHead {
		head(w, r)
		return
	} else if m == http.MethodGet {
		get(w, r)
		return
	} else if m == http.MethodPut {
		put(w, r)
		return
	} else if m == http.MethodPatch {
		patch(w, r)
		return
	} else if m == http.MethodPost {
		post(w, r)
		return
	} else if m == http.MethodDelete {
		del(w, r)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}













