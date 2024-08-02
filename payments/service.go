package main

import (
	"context"

	pb "github.com/swarajroy/oms-common/api"
)

type service struct {
}

func NewService() *service {
	return &service{}
}

func (s *service) CreatePayment(ctx context.Context, order *pb.Order) (string, error) {
	return "", nil
}
