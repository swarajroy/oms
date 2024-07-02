package main

import (
	"context"
	"log"

	common "github.com/swarajroy/oms-common"
	pb "github.com/swarajroy/oms-common/api"
)

type service struct {
	store OrderStore
}

func NewService(store OrderStore) *service {
	return &service{store: store}
}

func (s *service) CreateOrder(context.Context) error {
	return nil
}

func (s *service) ValidateOrder(ctx context.Context, order *pb.CreateOrderRequest) error {
	if len(order.Items) == 0 {
		return common.ErrNoItems
	}

	mergedItems := mergeItemsWithQuantity(order.Items)
	log.Println(mergedItems)
	return nil
}

func mergeItemsWithQuantity(itemsWithQuantity []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity {
	merged := make([]*pb.ItemsWithQuantity, 0)

	for _, itemWithQuantity := range itemsWithQuantity {
		found := false
		for _, finalItem := range merged {
			if finalItem.Id == itemWithQuantity.Id {
				finalItem.Quantity += itemWithQuantity.Quantity
				found = true
				break
			}
		}
		if !found {
			merged = append(merged, itemWithQuantity)
		}

	}
	return merged
}
