package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"encoding/json"
	"context"
)

type AckType int 

const (
	Ack AckType = iota
	NackRequeue
	NackDiscard
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	jsonData, err := json.Marshal(val)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{ContentType: "application/json", Body: jsonData}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, msg)
	if err != nil {
		return err
	}
	return nil
}

func SubscribeJSON[T any](
    conn *amqp.Connection,
    exchange,
    queueName,
    key string,
    queueType SimpleQueueType, // an enum to represent "durable" or "transient"
    handler func(T) AckType,
) error {
	err := subscribe(conn, exchange, queueName, key, queueType, handler, jsonUnmarshaller)
	if err != nil {
		return err
	}
	return nil
}

func jsonUnmarshaller[T any](data []byte) (T, error) {
	var msg T
	err := json.Unmarshal(data, &msg)
	return msg, err
}