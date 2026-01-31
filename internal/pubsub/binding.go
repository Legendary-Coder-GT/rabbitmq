package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int 

const (
	Durable SimpleQueueType = iota
	Transient
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return channel, amqp.Queue{}, err
	}

	var queue amqp.Queue

	if queueType == Durable {
		queue, err = channel.QueueDeclare(queueName, true, false, false, false, nil)
	} else {
		queue, err = channel.QueueDeclare(queueName, false, true, true, false, nil)
	}
	if err != nil {
		return channel, amqp.Queue{}, err
	}

	err = channel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return channel, amqp.Queue{}, err
	}

	return channel, queue, nil
}