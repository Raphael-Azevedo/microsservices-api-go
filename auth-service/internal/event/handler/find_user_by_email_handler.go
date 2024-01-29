package handler

import (
	"biz-hub-auth-service/pkg/events"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/streadway/amqp"
)

type FindUserByEmailHandler struct {
	RabbitMQChannel *amqp.Channel
}

func NewFindUserByEmailHandler(rabbitMQChannel *amqp.Channel) *FindUserByEmailHandler {
	return &FindUserByEmailHandler{
		RabbitMQChannel: rabbitMQChannel,
	}
}

func (h *FindUserByEmailHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("User Found: %v", event.GetPayload())
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
