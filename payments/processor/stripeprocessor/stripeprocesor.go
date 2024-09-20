package stripeprocessor

import (
	"context"
	"fmt"
	"log"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/checkout/session"
	common "github.com/swarajroy/oms-common"
	pb "github.com/swarajroy/oms-common/api"
)

var (
	GATEWAY_URL_ADDR = common.EnvString("GATEWAY_URL_ADDR", "http://localhost:8080")
)

type StripeProcessor struct {
}

func NewProcessor() *StripeProcessor {
	return &StripeProcessor{}
}

func (s *StripeProcessor) CreatePaymentLink(ctx context.Context, order *pb.Order) (string, error) {
	log.Printf("creating a paymentLink for order = %+v\n", order)

	gatewaySuccessURL := fmt.Sprintf("%s/success.html?customerId=%s&orderId=%s", GATEWAY_URL_ADDR, order.CustomerId, order.Id)
	gatewayCancelURL := fmt.Sprintf("%s/cancel.html", GATEWAY_URL_ADDR)

	items := []*stripe.CheckoutSessionLineItemParams{}
	for _, item := range order.Items {
		items = append(items, &stripe.CheckoutSessionLineItemParams{
			Price:    stripe.String(item.PriceID),
			Quantity: stripe.Int64(int64(item.Quantity)),
		})
	}

	params := &stripe.CheckoutSessionParams{
		Metadata: map[string]string{
			"orderId":    order.Id,
			"customerId": order.CustomerId,
		},
		LineItems:  items,
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(gatewaySuccessURL),
		CancelURL:  stripe.String(gatewayCancelURL),
	}

	result, err := session.New(params)

	if err != nil {
		return "", err
	}

	return result.URL, nil
}
