package rabbitmq

import (
	"api_service/common"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)
var logger *log.Logger
func init() {
	logger = common.WriteLog("logs/log.file")
}

func (q *RabbitMQ)Publish(exchange string, body interface{})  {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	e = q.Channel.Publish(exchange,
		"",
		false,
		false,
		amqp.Publishing{
		ReplyTo: q.Name,
		Body: []byte(str),
		})
	if e != nil {
		logger.Println(e)
		panic(e)
	}
}
