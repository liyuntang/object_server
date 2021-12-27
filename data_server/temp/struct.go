package temp

import (
	"data_server/locate"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type tempInfo struct {
	//Uuid string
	Name string
	Size int64
}

func (t *tempInfo)writeToFile() (name string, e error) {
	fileName := fmt.Sprintf("%s/temp/%s", locate.STORAGE_ROOT, t.Name)
	// 判断文件是否存在，如果已经存在则直接返回201
	base := strings.Split(fileName, ".")[0]

	list, _ := filepath.Glob(fmt.Sprintf("%s.*", base))
	//logger.Println(">>>>>>>>>>>>>", list)
	//logger.Println(fileName)
	if len(list) > 0 {
		for _, f := range list {
			fName := filepath.Base(strings.Trim(f, " "))
			if !strings.HasSuffix(fName, "dat") {
				return fName,nil
			}
		}
	}

	file, err := os.Create(fileName)
	if err != nil {
		//fmt.Println("create file", fileName, "is bad")
		return "", err
	}
	//fmt.Println("create file", fileName, "is ok")
	defer file.Close()
	buf, _ := json.Marshal(t)
	file.Write(buf)
	//logger.Println(fileName, "write data is ok")
	return "", nil
}

func (t *tempInfo)hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]
}

func (t *tempInfo)id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id
}



