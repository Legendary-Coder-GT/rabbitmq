package main

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	
	"fmt"
)

func handlerMove(gs *gamelogic.GameState) func(gamelogic.ArmyMove) {
	defer fmt.Print(">")
	return func(mv gamelogic.ArmyMove) {
		gs.HandleMove(mv)
	}
}