package gateway

import (
	"context"

	pb "github.com/swarajroy/oms-common/api"
)

type OrderGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order, error)
	GetOrder(ctx context.Context, customerId, orderId string) (*pb.Order, error)
}
