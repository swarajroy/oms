package main

import "context"

type OrderService interface {
	CreateOrder(context.Context) error
}

type OrderStore interface {
	Create(context.Context) error
}
