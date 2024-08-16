package main

import (
	"context"

	pb "github.com/swarajroy/oms-common/api"
)

type OrderService interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrder(context.Context, *pb.GetOrderRequest) (*pb.Order, error)
}

type OrderStore interface {
	Create(context.Context, *pb.Order) error
	Get(ctx context.Context, orderId string) (*pb.Order, error)
}
