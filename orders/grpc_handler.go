package main

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	pb "github.com/swarajroy/oms-common/api"
	"github.com/swarajroy/oms-common/broker"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	channel *amqp.Channel
	qName   string
	service *service
}

func NewGRPCHandler(grpcServer *grpc.Server, channel *amqp.Channel, service *service) {
	q, err := channel.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	handler := &grpcHandler{channel: channel, qName: q.Name, service: service}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New Order Received! Order = %v", p)
	o, err := h.service.CreateOrder(ctx, p)
	if err != nil {
		return nil, err
	}
	marshalledOrder, err := json.Marshal(o)
	if err != nil {
		return nil, err
	}

	err = h.channel.PublishWithContext(ctx, "", h.qName, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         marshalledOrder,
	})

	if err != nil {
		return nil, err
	}
	return o, nil
}

func (h *grpcHandler) GetOrder(ctx context.Context, p *pb.GetOrderRequest) (*pb.Order, error) {
	return h.service.GetOrder(ctx, p)
}
