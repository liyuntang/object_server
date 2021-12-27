package version

import (
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method

	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//from := 0
	//to := 1000
	//name := strings.Split(r.URL.EscapedPath(),"/")[2]
	// 说明要获取所有对象的所有版本，也就是说要获取对象列表，这个后边实现
}
