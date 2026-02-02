package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	
	"fmt"
	"time"
)

func handlerWar(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	defer fmt.Print(">")
	return func(rw gamelogic.RecognitionOfWar) pubsub.AckType {
		out, winner, loser := gs.HandleWar(rw)
		switch out {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon, gamelogic.WarOutcomeYouWon:
			msg := fmt.Sprintf("%s won a war against %s", winner, loser)
			return PublishGameLog(gs, msg)
		case gamelogic.WarOutcomeDraw:
			msg := fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			return PublishGameLog(gs, msg)
		default:
			fmt.Print("Error determining war outcome\n")
			return pubsub.NackDiscard
		}
	}
}

func PublishGameLog(gs *gamelogic.GameState, msg string) pubsub.AckType {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		fmt.Print("could not connect to RabbitMQ: %v", err)
		return pubsub.NackDiscard
	}
	defer conn.Close()
	
	channel, err := conn.Channel()

	exchange := routing.ExchangePerilTopic
	routingKey := routing.GameLogSlug + "." + gs.GetUsername()

	data := routing.GameLog{
		CurrentTime: time.Now(),
		Message: msg,
		Username: gs.GetUsername(),
	}

	err = pubsub.PublishGob(channel, exchange, routingKey, data)
	if err != nil {
		fmt.Print("error publishing move: %v", err)
		return pubsub.NackRequeue
	}
	return pubsub.Ack
}