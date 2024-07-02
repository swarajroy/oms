package main

import (
	"errors"
	"net/http"

	common "github.com/swarajroy/oms-common"
	pb "github.com/swarajroy/oms-common/api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type handler struct {
	client pb.OrderServiceClient
}

func NewHandler(client pb.OrderServiceClient) *handler {
	return &handler{client: client}
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

	if err := validateItems(items); err != nil {
		common.WriteError(w, err.Error(), http.StatusBadRequest)
		return
	}

	o, err := h.client.CreateOrder(r.Context(), &pb.CreateOrderRequest{
		CustomerId: r.PathValue("customerId"),
		Items:      items,
	})
	rStatus := status.Convert(err)
	if rStatus != nil {
		if rStatus.Code() != codes.OK {
			common.WriteError(w, rStatus.Message(), http.StatusBadRequest)
			return
		}
	}

	if err != nil {
		common.WriteError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	common.WriteJSON(w, o, http.StatusOK)

}

func validateItems(items []*pb.ItemsWithQuantity) error {
	if len(items) == 0 {
		return common.ErrNoItems
	}

	for _, item := range items {
		if len(item.Id) == 0 {
			return errors.New("item does not have a valid Id")
		}

		if item.Quantity <= 0 {
			return errors.New("item does not have a valid Quantity")
		}
	}

	return nil
}
