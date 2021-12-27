package rabbitmq

import (
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	channel *amqp.Channel
	Name string
	exchange string
}

func New(s string) *RabbitMQ {
	conn, e := amqp.Dial(s)
	if e != nil {
		panic(e)
	}
	ch, e := conn.Channel()
	if e != nil {
		panic(e)
	}
	q, e := ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil,)

	if e != nil {
		panic(e)
	}
	mq := new(RabbitMQ)
	mq.channel = ch
	mq.Name = q.Name
	return mq
}

func (q *RabbitMQ)Close()  {
	q.channel.Close()
}