package main

import (
	"fmt"

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

	gs := gamelogic.NewGameState(name)
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, routing.PauseKey+"."+name, routing.PauseKey, "Transient", handlerPause(gs))
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}

	gameloop(gs)
}
