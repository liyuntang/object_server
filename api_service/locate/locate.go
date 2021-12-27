package locate

import (
	"api_service/common"
	"api_service/heartbeat"
	"api_service/rabbitmq"
	"api_service/rs"
	"api_service/types"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"strings"
	"time"
)
var logger *log.Logger

func init() {
	logger = common.WriteLog("logs/log.file")
}

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[2])
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(info)
	w.Write(b)
}

func Locate(name string) (locateInfo map[int]string) {
	conn, e := amqp.Dial(heartbeat.RABBITMQ_SERVER)
	if e != nil {
		logger.Println(e)
		panic(e)
	}
	defer conn.Close()
	ch, e := conn.Channel()
	if e != nil {
		logger.Println(e)
		panic(e)
	}
	defer ch.Close()
	qu, e := ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,)

	if e != nil {
		logger.Println(e)
		panic(e)
	}
	q := &rabbitmq.RabbitMQ{}
	q.Name = qu.Name
	q.Channel = ch
	q.Publish("dataServers", name)
	c := q.Consume()
	okCh := make(chan int)
	locateInfo = make(map[int]string)
	go func() {
		select {
		case <-time.After(time.Second):
			return
		case <-okCh:
			return
		}
	}()
	for i:=0;i< rs.ALL_SHARDS;i++ {
		msg := <- c
		if len(msg.Body) == 0 {
			return
		}
		var info types.LocateMessage
		json.Unmarshal(msg.Body, &info)
		if info.Id != -1 {
			locateInfo[info.Id] = info.Addr
		}
	}

	okCh <- 1
	close(okCh)
	return
}

func Exist(name string) bool {
	locateInfo := Locate(name)
	return len(locateInfo) >= rs.DATA_SHARDS
}