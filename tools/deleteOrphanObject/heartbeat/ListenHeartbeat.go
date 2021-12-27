package heartbeat

import (
	"deleteOrphanObject/rabbitmq"
	"flag"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var apiServers = make(map[string]time.Time)
var mutex sync.Mutex
var RABBITMQ_SERVER string

func init() {
	flag.StringVar(&RABBITMQ_SERVER, "r", "amqp://admin:123.com@10.10.10.62:5672", "rabbit server")
}

func ListenHeartbeat() {
	q := rabbitmq.New(RABBITMQ_SERVER)
	defer q.Close()
	q.Bind("online")
	c := q.Consume()
	go removeExpiredDataServer()
	for msg := range c {
		apiServer, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		//fmt.Println(time.Now().Format("2006-01-02 15:04:05"), ">>>>>>>>>>>>", apiServer)
		//fmt.Println(".............................")
		apiServers[apiServer] = time.Now()
		mutex.Unlock()
	}
}

func removeExpiredDataServer() {
	for {
		time.Sleep(1*time.Second)
		mutex.Lock()
		for s, t := range apiServers {
			if t.Add(3*time.Second).Before(time.Now()) {
				delete(apiServers, s)
			}
		}
		mutex.Unlock()
		//fmt.Println(">>>>>>>>>>>>", dataServers)
	}
}

func getApiServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range apiServers {
		ds = append(ds, s)
	}
	return ds
}

func ChooseRandomDataServer() string {
	servers := getApiServers()
	//fmt.Println(servers)
	num := len(servers)
	if num == 0 {
		return ""
	}
	s := servers[rand.Intn(num)]
	return s
}





