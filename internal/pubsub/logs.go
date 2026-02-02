package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"encoding/gob"
	"context"
	"bytes"
)

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(val)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{ContentType: "application/gob", Body: b.Bytes()}
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, msg)
	if err != nil {
		return err
	}
	return nil
}

func SubscribeGob[T any](
    conn *amqp.Connection,
    exchange,
    queueName,
    key string,
    queueType SimpleQueueType, // an enum to represent "durable" or "transient"
    handler func(T) AckType,
) error {
	err := subscribe(conn, exchange, queueName, key, queueType, handler, gobUnmarshaller)
	if err != nil {
		return err
	}
	return nil
}

func gobUnmarshaller[T any](data []byte) (T, error) {
	b := bytes.NewBuffer(data)
	dec := gob.NewDecoder(b)
	var gl T
	err := dec.Decode(&gl)
	if err != nil {
		return gl, err
	}
	return gl, nil
}