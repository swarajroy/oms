package inmem

import (
	"context"

	pb "github.com/swarajroy/oms-common/api"
)

type InMemProcessor struct {
}

func NewInMemProcessor() *InMemProcessor {
	return &InMemProcessor{}
}

func (in *InMemProcessor) CreatePaymentLink(ctx context.Context, order *pb.Order) (string, error) {
	return "dummy-link", nil
}
