package main

import (
	"context"

	pb "github.com/swarajroy/oms-common/api"
)

type OrderService interface {
	CreateOrder(context.Context) error
	ValidateOrder(context.Context, *pb.CreateOrderRequest)
}

type OrderStore interface {
	Create(context.Context) error
}
