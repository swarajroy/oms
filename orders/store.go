package main

import (
	"context"

	pb "github.com/swarajroy/oms-common/api"
)

var orders = make(map[string]*pb.Order)

type store struct {
}

func NewStore() *store {
	return &store{}
}

func (s *store) Create(ctx context.Context, p *pb.Order) error {
	orders[p.Id] = p
	return nil
}

func (s *store) Get(ctx context.Context, orderId string) (*pb.Order, error) {
	return orders[orderId], nil
}
