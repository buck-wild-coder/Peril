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

	gamelogic.PrintServerHelp()
	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		if words[0] == "pause" {
			fmt.Println("sending a pause message")
		} else if words[0] == "resume" {
			fmt.Println("sending a resume message")
		} else if words[0] == "quit" {
			break
		} else {

		}
	}

	Exit()
}

func Exit() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("\nThe program is shutting down and close the connection.")
}
