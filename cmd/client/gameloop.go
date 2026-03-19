package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func gameloop(gs *gamelogic.GameState) {
	var err error
	for {
		words := gamelogic.GetInput()
		switch words[0] {
		case "spawn":
			err = gs.CommandSpawn(words)
			if err != nil {
				fmt.Println("msg", err)
				continue
			}

		case "move":
			_, err = gs.CommandMove(words)
			if err != nil {
				fmt.Println("msg", err)
				continue
			}
			fmt.Println("it worked.")

		case "status":
			gs.CommandStatus()

		case "help":
			gamelogic.PrintClientHelp()

		case "spam":
			fmt.Println("Spamming not allowed yet!")

		case "quit":
			gamelogic.PrintQuit()
			Exit()

		default:
			fmt.Println("error message: Unknown Command")
		}
	}
}

func Exit() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nThe program is shutting down and close the connection.")
}

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) {
	return func(ps routing.PlayingState) {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
	}
}
