package main

import (
	"context"
	"testing"

	"github.com/swarajroy/oms-common/api"
	"github.com/swarajroy/oms-payments/processor/stripeprocessor/inmem"
)

func TestService(t *testing.T) {
	processor := inmem.NewInMemProcessor()
	svc := NewService(processor)

	t.Run("should create a payment link", func(t *testing.T) {
		link, err := svc.CreatePayment(context.Background(), &api.Order{})

		if err != nil {
			t.Error("service create payment failed")
		}

		if link == "" {
			t.Error("service create payment failed")
		}
	})

}
