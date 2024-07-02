package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"

	pb "github.com/swarajroy/oms-common/api"
	"google.golang.org/grpc"
)

type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service *service
}

func NewGRPCHandler(grpcServer *grpc.Server, service *service) {
	handler := &grpcHandler{service: service}
	pb.RegisterOrderServiceServer(grpcServer, handler)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, p *pb.CreateOrderRequest) (*pb.Order, error) {

	log.Printf("New Order Received! Order = %v", p)

	o := &pb.Order{
		Id:         strconv.Itoa(rand.Int()),
		CustomerId: p.CustomerId,
		Status:     "PENDING",
	}
	return o, nil
}
