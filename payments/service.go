package main

import (
	"context"

	pb "github.com/swarajroy/oms-common/api"
	"github.com/swarajroy/oms-payments/processor"
)

type service struct {
	processor processor.PaymentProcessor
}

func NewService(processor processor.PaymentProcessor) *service {
	return &service{processor: processor}
}

func (s *service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	paymentLink, err := s.processor.CreatePaymentLink(ctx, order)
	if err != nil {
		return "", err
	}
	return paymentLink, nil
}
