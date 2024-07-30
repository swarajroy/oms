package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/swarajroy/oms-common/api"
	"github.com/swarajroy/oms-common/broker"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	channel *amqp.Channel
	qName   string
}

func NewGRPCHandler(grpcServer *grpc.Server, channel *amqp.Channel) {
	q, err := channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	handler := &grpcHandler{channel: channel, qName: q.Name}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New Order Received! Order = %v", p)
	o := &pb.Order{
		Id:         strconv.Itoa(rand.Int()),
		CustomerId: p.CustomerId,
		Status:     "PENDING",
	}
	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}
	
	h.channel.PublishWithContext(ctx, "", h.qName, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         marshalledOrder,
	})
	return o, nil
}
