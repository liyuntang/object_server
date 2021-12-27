package temp

import (
	"data_server/locate"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)



func post(w http.ResponseWriter, r *http.Request) {
	//logger.Println(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> post")
	name := strings.Split(r.URL.EscapedPath(),"/")[2]
	size, err := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t := tempInfo{name, size}
	file, err := t.writeToFile()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(file) != 0 {
		// 说明该对象已经创建过了，可能是断点续传操作，此时直接返回201
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(file))
		return
	}
	// 这里是根据uuid创建的文件，所以每个dataServer上创建的文件是不一样的，同时每个data server生成的uuid返回给了api server并且放入token中
	fileName := fmt.Sprintf("%s/temp/%s.dat", locate.STORAGE_ROOT, t.Name)
	_, err = os.Create(fileName)
	if err != nil {
		//logger.Println("create file", fileName, "is bad")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//logger.Println("create file", fileName, "is ok")
	w.WriteHeader(http.StatusCreated)
	return
}







