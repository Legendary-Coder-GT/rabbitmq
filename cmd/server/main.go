package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"os"
	"os/signal"
	"log"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!\n")

	gamelogic.PrintServerHelp()

	channel, err := conn.Channel()

	exchange := routing.ExchangePerilTopic
	queueName := routing.GameLogSlug
	routingKey := routing.GameLogSlug + ".*"
	queueType := pubsub.Durable

	_, _, err = pubsub.DeclareAndBind(conn, exchange, queueName, routingKey, queueType)
	if err != nil {
		log.Fatalf("error declaring and binding queue: %v", err)
		return
	}

	for ;; {
		input := gamelogic.GetInput()
		cmd := input[0]

		if cmd == "pause" {
			fmt.Println("Sending pause message...\n")
			data := routing.PlayingState{true}
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, data)

		} else if cmd == "resume" {
			fmt.Println("Sending resume message...\n")
			data := routing.PlayingState{false}
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, data)

		} else if cmd == "quit" {
			fmt.Println("Quitting...\n")
			break
		} else {
			fmt.Println("I don't understand your command\n")
		}
	}
}
