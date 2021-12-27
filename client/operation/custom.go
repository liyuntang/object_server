package operation

import (
	"flag"
	"os"
	"path/filepath"
)

var tmpDir string
//var bigOjbectSize int64
const bigOjbectSize = 1024*1024*5
func init() {
	flag.StringVar(&tmpDir, "tmp", "/Users/liyuntang/git/object_server/client/tmp", "临时目录,需要使用绝对路径")

}

func Custom(object string) string{
	info, err := os.Stat(object)
	if err != nil {
		logger.Println("获取", object, "信息失败")
		return "0"
	}
	oInfo := &objectInfo{}
	oInfo.name = filepath.Base(object)
	oInfo.size = info.Size()
	oInfo.path = object
	if oInfo.size > bigOjbectSize {
		oInfo.isBigObject = true
	}
	return oInfo.doit()
}
