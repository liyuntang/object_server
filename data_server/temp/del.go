package temp

import (
	"data_server/locate"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func del(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	infoFile := fmt.Sprintf("%s/temp/%s.dat", locate.STORAGE_ROOT, uuid)
	datFile := fmt.Sprintf("%s.dat", infoFile)
	os.Remove(infoFile)
	os.Remove(datFile)
}

