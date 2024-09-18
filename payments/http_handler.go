package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
	common "github.com/swarajroy/oms-common"
)

type PaymentHttpHandler struct {
}

func NewPaymentHttpHandler() *PaymentHttpHandler {
	return &PaymentHttpHandler{}
}

func (ph *PaymentHttpHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/webhook", ph.HandleSomeRoute)
}

func (ph *PaymentHttpHandler) HandleSomeRoute(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse webhook body json: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// If you are testing your webhook locally with the Stripe CLI you
	// can find the endpoint's secret by running `stripe listen`
	// Otherwise, find your endpoint's secret in your webhook settings
	// in the Developer Dashboard
	endpointSecret := common.EnvString("ENDPOINT_SECRET", "")

	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err = webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Then define and call a func to handle the successful payment intent.
		// handlePaymentIntentSucceeded(paymentIntent)
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Then define and call a func to handle the successful attachment of a PaymentMethod.
		// handlePaymentMethodAttached(paymentMethod)
	// ... handle other event types
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if session.PaymentStatus == "paid" {
			log.Printf("Payment for checkout session %v succeeded", session.ID)
		}
		
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}
