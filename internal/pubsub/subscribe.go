package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func subscribe[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType,
	handler func(T) AckType,
	unmarshaller func([]byte) (T, error),
) error {
	channel, queue, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	err = channel.Qos(10, 0, true)
	if err != nil {
		return err
	}

	ch, err := channel.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for delivery := range ch {
			msg, err := unmarshaller(delivery.Body)
			if err != nil {
				log.Print(err)
				delivery.Ack(false)
				continue
			}
			msg_type := handler(msg)
			switch msg_type {
			case Ack:
				delivery.Ack(false)
				log.Print("Acknowledged\n")
				break
			case NackRequeue:
				delivery.Nack(false, true)
				log.Print("Negative Acknolwedged and Requeued\n")
				break
			case NackDiscard:
				delivery.Nack(false, false)
				log.Print("Negative Acknolwedged and Discarded\n")
				break
			}
		}
	}()

	return nil

}