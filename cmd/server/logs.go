package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	
	"fmt"
)

func handlerLogs() func(routing.GameLog) pubsub.AckType {
	defer fmt.Print(">")
	return func(gl routing.GameLog) pubsub.AckType {
		gamelogic.WriteLog(gl)
		return pubsub.Ack
	}
}