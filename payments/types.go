package main

import (
	"context"

	pb "github.com/swarajroy/oms-common/api"
)

type PaymentsService interface {
	CreatePayment(context.Context, *pb.Order) (string, error)
}
