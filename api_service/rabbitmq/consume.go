package rabbitmq

import "github.com/streadway/amqp"

func (q *RabbitMQ)Consume() <-chan amqp.Delivery{
	//logger.Println(">>>>>>>>>>>>>> consume")
	c, e := q.Channel.Consume(q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,)

	if e != nil {
		panic(e)
	}
	return c
}
