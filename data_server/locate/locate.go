package locate

import (
	"data_server/common"
	"data_server/heartbeat"
	"data_server/monit"
	"data_server/rabbitmq"
	"data_server/types"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)
var objects = make(map[string]int)
var mutex sync.RWMutex
var STORAGE_ROOT string
var logger *log.Logger
var hostName, _ = os.Hostname()
func init() {
	flag.StringVar(&STORAGE_ROOT, "storage", "", "storage dir")
	logger = common.WriteLog("logs/log.file")
}
func Locate(hash string) int {
	mutex.RLock()
	id, ok := objects[hash]
	mutex.RUnlock()
	// ok=true表示对象存在 ok=false表示对象不存在
	if !ok {
		// 说明需要定位的对象不存在
		return -1
	}
	// 说明对象存在
	return id
}

func Add(hash string, id int) {
	mutex.Lock()
	objects[hash] = id
	mutex.Unlock()
	monit.Object_total.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS).Inc()
}

func Del(hash string)  {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
	monit.Object_total.WithLabelValues(hostName, heartbeat.LISTEN_ADDRESS).Dec()
}
func StartLocate() {
	q := rabbitmq.New(heartbeat.RABBITMQ_SERVER)
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		// strconv.Unquote的作用是将输入的字符串的双引号去掉，并返回去掉双引号的内容
		hash, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		//bTime := time.Now()
		id := Locate(hash)
		// id = -1说明定位的对象不存在	id!=-1表示对象存在
		info := types.LocateMessage{}
		info.Addr = heartbeat.LISTEN_ADDRESS
		info.Id = id
		if id == -1 {
			info.ISExsit = false
		} else {
			// id!=-1表示对象存在
			info.ISExsit = true
		}
		q.Send(msg.ReplyTo, info)
	}
}

func CollectObjects() {
	//logger.Println("start collect objects")
	files, _ := filepath.Glob(fmt.Sprintf("%s/objects/*", STORAGE_ROOT))
	for i := range files {
		file := strings.Split(filepath.Base(files[i]), ".")
		if len(file) != 3 {
			panic(files[i])
		}
		hash := file[0]
		id, e := strconv.Atoi(file[1])
		//logger.Println("id is", id)
		if e != nil {
			panic(e)
		}
		//objects[hash] = id
		Add(hash, id)
	}
	//logger.Println("collect object is over, objects is", objects)
}
