package handler

import (
	"biz-hub-auth-service/pkg/events"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)

type UserCreatedHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewUserCreatedHandler(rabbitMQChannel *amqp.Channel) *UserCreatedHandler {
	return &UserCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *UserCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("User created: %v", event.GetPayload())
	jsonOutput, _ := json.Marshal(event.GetPayload())

	msgRabbitmq := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonOutput,
	}

	h.RabbitMQChannel.Publish(
		"amq.direct", // exchange
		"",           // key name
		false,        // mandatory
		false,        // immediate
		msgRabbitmq,  // message to publish
	)
}
