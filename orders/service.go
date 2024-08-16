package main

import (
	"context"

	pb "github.com/swarajroy/oms-common/api"
)

type service struct {
	store OrderStore
}

func NewService(store OrderStore) *service {
	return &service{store: store}
}

func (s *service) CreateOrder(ctx context.Context, order *pb.CreateOrderRequest) (*pb.Order, error) {
	items := []*pb.Item{}
	for _, itemRequest := range order.Items {
		item := &pb.Item{
			Id:       itemRequest.Id,
			Quantity: itemRequest.Quantity,
			Name:     "cheese burger",
			PriceID:  "price_1PlBysFmQPISRdloYnGomtYm",
		}
		items = append(items, item)
	}

	o := &pb.Order{
		Id:         "42",
		CustomerId: order.CustomerId,
		Status:     "pending",
		Items:      items,
	}
	err := s.store.Create(ctx, o)

	if err != nil {
		return nil, err
	}

	return o, nil
}

func (s *service) GetOrder(ctx context.Context, order *pb.GetOrderRequest) (*pb.Order, error) {
	o, err := s.store.Get(ctx, order.OrderId)
	if err != nil {
		return nil, err
	}
	return o, nil
}
