package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	data, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("Can't marshal: %v", err)
	}

	ctx := context.Background()
	err = ch.PublishWithContext(ctx, exchange, key, false, false, amqp.Publishing{ContentType: "application/json", Body: data})
	if err != nil {
		return fmt.Errorf("Can't publish the error..%v", err)
	}

	return nil
}

type SimpleQueueType string

const (
	Durable   SimpleQueueType = "Durable"
	Transient SimpleQueueType = "Transient"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // SimpleQueueType is an "enum" type I made to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("Can't create channel: %v", err)
	}

	amqpQueue, err := ch.QueueDeclare(queueName, queueType == "Durable", queueType == "Transient", queueType == "Transient", false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("Can't declare the queue: %v", err)
	}

	err = ch.QueueBind(amqpQueue.Name, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("Can't bind the queue: %v", err)
	}

	return ch, amqpQueue, nil
}
