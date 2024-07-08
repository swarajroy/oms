package main

import (
	"net/http"

	common "github.com/swarajroy/oms-common"
	pb "github.com/swarajroy/oms-common/api"
	"github.com/swarajroy/oms-gateway/gateway"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	orderGateway gateway.OrderGateway
}

func NewHandler(orderGateway gateway.OrderGateway) *handler {
	return &handler{orderGateway: orderGateway}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/1/customers/{customerId}/orders", h.HandleCreateOrder)
}

func (h *handler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {

	var items []*pb.ItemsWithQuantity
	if err := common.ReadJSON(r, &items); err != nil {
		common.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	o, err := h.orderGateway.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerId: r.PathValue("customerId"),
		Items:      items,
	})

	rStatus := status.Convert(err)
	
	if rStatus.Code() != codes.OK {
		common.WriteError(w, rStatus.Message(), http.StatusBadRequest)
		return
	}

	common.WriteJSON(w, o, http.StatusOK)
}
