package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	
	"fmt"
)

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType {
	defer fmt.Print(">")
	return func(mv gamelogic.ArmyMove) pubsub.AckType {
		out := gs.HandleMove(mv)
		switch out {
		case gamelogic.MoveOutComeSafe:
			return pubsub.Ack
		case gamelogic.MoveOutcomeMakeWar:
			const rabbitConnString = "amqp://guest:guest@localhost:5672/"

			conn, err := amqp.Dial(rabbitConnString)
			if err != nil {
				fmt.Print("could not connect to RabbitMQ: %v", err)
				return pubsub.NackDiscard
			}
			defer conn.Close()
			
			channel, err := conn.Channel()

			exchange := routing.ExchangePerilTopic
			queueName := routing.WarRecognitionsPrefix + "." + gs.GetUsername()

			data := gamelogic.RecognitionOfWar{
				Attacker: mv.Player,
				Defender: gs.GetPlayerSnap(),
			}

			err = pubsub.PublishJSON(channel, exchange, queueName, data)
			if err != nil {
				fmt.Print("error publishing move: %v", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			return pubsub.NackDiscard
		}
	}
}