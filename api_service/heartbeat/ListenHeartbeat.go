package heartbeat

import (
	"api_service/common"
	"api_service/rabbitmq"
	"encoding/json"
	"flag"
	"github.com/shirou/gopsutil/disk"
	"log"
	"math/rand"
	"sync"
	"time"
)

var logger *log.Logger

var dataServers = make(map[DataServerInfo]time.Time)
var mutex sync.Mutex
var RABBITMQ_SERVER string
var LISTEN_ADDRESS string

/*
这个地方将根据data server几点的磁盘剩余容量进行权重配置，权重计算规则：
	1、weight := 100 - info.UsageStat.UsedPercent
	2、如果weight >= 90则不在分配
	3、weight越小分配的几率越高
*/
//
var dataServerChannel = make(chan string, 1000)


func init() {
	flag.StringVar(&RABBITMQ_SERVER, "r", "amqp://admin:123.com@10.10.30.207:5672", "rabbit server")
	flag.StringVar(&LISTEN_ADDRESS, "listen", "", "api server")
	logger = common.WriteLog("logs/log.file")
}
type DataServerInfo struct {
	HostName string
	Address string
	UsageStat *disk.UsageStat
}

func ListenHeartbeat() {
	//logger.Println("connect to rabbitmq")
	q := rabbitmq.New(RABBITMQ_SERVER)
	defer q.Close()
	q.Bind("apiServers")
	c := q.Consume()
	go removeExpiredDataServer()
	go func() {
		for {
			if len(dataServers) >0 {
				//for info, _ := range dataServers {
				//	weight := int(100 - info.UsageStat.UsedPercent)
				//	if weight < 90 {
				//		for i:=1;i<=weight;i++ {
				//			logger.Println(info.Address, weight)
				//			dataServerChannel <- info.Address
				//		}
				//	}
				//}
				logger.Println(">>>>>>>>>>", dataServers)
			} else {
				logger.Println("data servers is null")
				time.Sleep(time.Second)
			}
		}


	}()
	for msg := range c {
		dataServerInfo := &DataServerInfo{}
		e := json.Unmarshal(msg.Body, dataServerInfo)
		//dataServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		//logger.Println(">>>>>>>>>>>>>>", dataServerInfo)
		mutex.Lock()
		dataServers[*dataServerInfo] = time.Now()
		mutex.Unlock()
		//logger.Println(dataServers)
	}
	q.Close()
}

func removeExpiredDataServer() {
	for {
		time.Sleep(1*time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10*time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}




func ChooseRandomDataServers(n int, exclude map[int]string) (ds []string) {
	return []string{}
	candidates := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	for id, addr := range exclude {
		reverseExcludeMap[addr] = id
	}
	logger.Println("reverseExcludeMap is", reverseExcludeMap)

	for server := range dataServerChannel {
		//logger.Println("server is", server)
		_, excluded := reverseExcludeMap[server]
		if !excluded {
			// 说明不在排除的序列
			for _, d := range candidates {
				if d != server {
					candidates = append(candidates, server)
				}
			}


		}
	}
	length := len(candidates)
	if length < n {
		return
	}
	p := rand.Perm(length)
	for i:=0;i<n;i++ {
		ds = append(ds, candidates[p[i]])
	}
	logger.Println(">>>>>>>>>", candidates)
	return
}

// 向online注册自己
func Regist() {
	//logger.Println("connect to rabbitmq")
	q := rabbitmq.New(RABBITMQ_SERVER)
	defer q.Close()
	for {
		q.Publish("online", LISTEN_ADDRESS)
		time.Sleep(1*time.Second)
	}
}





