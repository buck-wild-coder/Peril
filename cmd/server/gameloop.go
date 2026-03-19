package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	amqp "github.com/rabbitmq/amqp091-go"
)

func gameloop(amqpChann *amqp.Channel) {
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
