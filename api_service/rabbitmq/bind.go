package rabbitmq

func (q *RabbitMQ)Bind(exchange string)  {
	e := q.Channel.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil)

	if e != nil {
		panic(e)
	}

	q.exchange = exchange





}
