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

	uname, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("error logging in: %v", err)
		return
	}

	exchange := routing.ExchangePerilDirect
	queueName := routing.PauseKey + "." + uname
	routingKey := routing.PauseKey
	queueType := pubsub.Transient

	_, _, err = pubsub.DeclareAndBind(conn, exchange, queueName, routingKey, queueType)
	if err != nil {
		log.Fatalf("error declaring and binding queue: %v", err)
		return
	}

	gs := gamelogic.NewGameState(uname)
	gamelogic.PrintClientHelp()

	for ;; {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "move":
			_, err := gs.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
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
