package main

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/swarajroy/oms-common/api"
	"github.com/swarajroy/oms-common/broker"
)

type consumer struct {
	service PaymentsService
	ch      *amqp.Channel
	q       amqp.Queue
}

func NewConsumer(svc PaymentsService, channel *amqp.Channel) *consumer {
	q, err := channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}
	return &consumer{service: svc, ch: channel, q: q}
}

func (c *consumer) Listen() {
	msgs, err := c.ch.Consume(c.q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	var forever chan struct{}
	go func() {
		for d := range msgs {
			o := &pb.Order{}
			err = json.Unmarshal(d.Body, o)
			if err != nil {
				log.Printf("Error occurred whilst unmarshalling order")
				continue
			}
			log.Printf("Received message :%+v\n", o)
			paymentLink, err := c.service.CreatePayment(context.Background(), o)
			if err != nil {
				log.Printf("Error occurred whilst creating the payment link for order %+v\n", o)
				continue
			}
			log.Printf("payment link = %s", paymentLink)
		}
	}()

	<-forever
}
