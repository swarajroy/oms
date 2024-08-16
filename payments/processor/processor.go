package processor

import (
	"context"

	pb "github.com/swarajroy/oms-common/api"
)

type PaymentProcessor interface {
	CreatePaymentLink(context.Context, *pb.Order) (string, error)
}
