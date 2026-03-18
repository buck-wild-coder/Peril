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
	fmt.Println("Starting Peril server...")
	url := "amqp://guest:guest@localhost:5672/"

	amqpConn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println("error...", err)
		return
	}
	defer amqpConn.Close()
	fmt.Println("Connection was successful")

	amqpChann, err := amqpConn.Channel()
	if err != nil {
		fmt.Println("Error occured..", err)
		return
	}
	defer amqpChann.Close()

	gamelogic.PrintServerHelp()
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "pause":
			fmt.Println("sending a pause message")
			publish(amqpChann, true)
		case "resume":
			fmt.Println("sending a resume message")
			publish(amqpChann, false)
		case "quit":
			fmt.Println("Quiting")
			Exit()
			return
		default:
			fmt.Println("don't understand the command.")
		}
	}
}

func publish(amqpChann *amqp.Channel, IsPaused bool) {
	err := pubsub.PublishJSON(amqpChann, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: IsPaused})
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func Exit() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nThe program is shutting down and close the connection.")
}
