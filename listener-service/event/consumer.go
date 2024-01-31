package event

import (
	"encoding/json"
	"fmt"
	"log"
	"net/rpc"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn      *amqp.Connection
	queueName string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return consumer, nil
}

func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	queue, err := declareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		ch.QueueBind(
			queue.Name,
			s,
			"logs_topic",
			false,
			nil,
		)
		if err != nil {
			return err
		}
	}
	messages, err := ch.Consume(queue.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for message := range messages {
			var payload Payload
			_ = json.Unmarshal(message.Body, &payload)

			go handlePayload(payload)
		}
	}()

	fmt.Printf("Waiting for message [Exchange, Queue] [logs_topic, %s]\n", queue.Name)
	<-forever

	return nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event", "authentication":
		err := logItemViaRPC(payload)
		if err != nil {
			log.Println(err)
		}

	case "auth":
		// authenticate
	case "logger":
		// add logger
	default:
		err := logItemViaRPC(payload)
		if err != nil {
			log.Println(err)
		}

	}
}

func logItemViaRPC(payload Payload) error {
	client, err := rpc.Dial("tcp", "logger-service:5001")
	if err != nil {
		return err
	}

	var result string
	err = client.Call("RPCServer.LogInfo", payload, &result)
	if err != nil {
		return err
	}
	log.Printf("Messager send by: %s", payload.Name)
	return nil
}
