package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	name, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	url := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	defer conn.Close()
	fmt.Println("Connection was successful")

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+name, routing.PauseKey, "Transient")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	gamestate := gamelogic.NewGameState(name)
	for {
		words := gamelogic.GetInput()
		switch words[0] {
		case "spawn":
			err = gamestate.CommandSpawn(words)
			if err != nil {
				fmt.Println("msg", err)
				continue
			}

		case "move":
			_, err = gamestate.CommandMove(words)
			if err != nil {
				fmt.Println("msg", err)
				continue
			}
			fmt.Println("it worked.")

		case "status":
			gamestate.CommandStatus()

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
