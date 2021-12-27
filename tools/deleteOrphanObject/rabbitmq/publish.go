package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

func (q *RabbitMQ)Publish(exchange string, body interface{})  {
	str, e := json.Marshal(body)
	if e != nil {
		panic(e)
	}
	e = q.channel.Publish(exchange,
		"",
		false,
		false,
		amqp.Publishing{
		ReplyTo: q.Name,
		Body: []byte(str),
		})
	if e != nil {
		panic(nil)
	}
}
