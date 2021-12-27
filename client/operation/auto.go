package operation

import (
	"client/common"
	"flag"
	"fmt"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

/*
bigWrite是用来上传大对象的，流程如下：
	1、发起post请求获取token
	2、根据上一步获取的token发起head请求查看当前该object已经上传多少数据了
	3、根据上一步获取到的进度发起put请求上传数据

*/

var logger *log.Logger
var objectSize int64
var threads int
var wg sync.WaitGroup

func init() {

	logger = common.WriteLog("logs/log.file")
	flag.Int64Var(&objectSize, "size", 1, "对象大小,单位MB，")
	flag.IntVar(&threads, "threads", 1, "线程数")

}

func Auto(db *gorm.DB) {
	check()
	logger.Println("进入自动测试模式")
	hostName, _ := os.Hostname()
	size := objectSize*1024*1024
	//根据size生成一个临时字符串
	cTime := time.Now()
	caseOfData := getData(size)
	if len(*caseOfData) == 0 {
		logger.Println("sorry, make case of data is bad, run time is", time.Since(cTime))
		os.Exit(0)
	}
	logger.Println("make case of data is ok, run time is", time.Since(cTime), "data size is", objectSize,"MB")
	for i := 1; i <= threads; i++ {
		wg.Add(1)
		go func(id int, hostName string) {
			defer wg.Done()
			data := Object_server_markbench{}
			data.Host_name = hostName
			data.Thread_id = id
			for {
				object := fmt.Sprintf("%s/%s-%d-%d-%d", tmpDir, hostName, id, objectSize, time.Now().UnixNano())
				sTime := time.Now()
				logger.Println("make data")
				err := createFile(object, caseOfData)
				if err != nil {
					logger.Println("createFile is bad, err is", err, "run time is", time.Since(sTime))
					return
				}
				logger.Println("createFile is ok, run time is", time.Since(sTime), "data size is", objectSize,"MB")
				// 写入数据
				dTime := time.Now()
				data.Code,_ = strconv.Atoi(Custom(object))
				data.Run_time = time.Since(dTime).Milliseconds()
				//db.Create(data)
				logger.Println("remove object")
				os.Remove(object)
				//os.Exit(0)
			}

		}(i, hostName)
		wg.Wait()
	}
}
