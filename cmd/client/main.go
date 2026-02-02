package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"log"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
		return
	}
	defer conn.Close()
	fmt.Println("Starting Peril client...")

	channel, err := conn.Channel()

	uname, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("error logging in: %v", err)
		return
	}

	gs := gamelogic.NewGameState(uname)
	gamelogic.PrintClientHelp()

	exchange := routing.ExchangePerilDirect
	queueName := routing.PauseKey + "." + uname
	routingKey := routing.PauseKey
	queueType := pubsub.Transient

	err = pubsub.SubscribeJSON(conn, exchange, queueName, routingKey, queueType, handlerPause(gs))
	if err != nil {
		log.Fatalf("error subscribing to pause queue: %v", err)
		return
	}

	exchange = routing.ExchangePerilTopic
	queueName = routing.ArmyMovesPrefix + "." + uname
	routingKey = routing.ArmyMovesPrefix + ".*"
	queueType = pubsub.Transient

	err = pubsub.SubscribeJSON(conn, exchange, queueName, routingKey, queueType, handlerMove(gs))
	if err != nil {
		log.Fatalf("error subscribing to army_moves queue: %v", err)
		return
	}

	err = pubsub.SubscribeJSON(conn, exchange, "war", routing.WarRecognitionsPrefix + ".*", pubsub.Durable, handlerWar(gs))
	if err != nil {
		log.Fatalf("error subscribing to war queue: %v", err)
		return
	}

	for ;; {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "move":
			mv, err := gs.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = pubsub.PublishJSON(channel, exchange, queueName, mv)
			if err != nil {
				log.Fatalf("error publishing move: %v", err)
				return
			}
		case "spawn":
			err = gs.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("unknown command")
		}
	}
}
