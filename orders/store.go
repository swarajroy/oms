package main

import (
	"context"
	"fmt"

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
	result := orders[orderId]
	if result == nil {
		return nil, fmt.Errorf("order with orderId %s not found", orderId)
	}
	return orders[orderId], nil
}
