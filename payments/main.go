package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/stripe/stripe-go/v79"
	common "github.com/swarajroy/oms-common"
	"github.com/swarajroy/oms-common/broker"
	"github.com/swarajroy/oms-common/discovery"
	"github.com/swarajroy/oms-common/discovery/consul"
	"github.com/swarajroy/oms-payments/processor/stripeprocessor"
	"google.golang.org/grpc"
)

var (
	SERVICE_NAME = "payments"
	GRPC_ADDR    = common.EnvString("GRPC_ADDR", "localhost:2002")
	CONSUL_ADDR  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	AMQP_USER    = common.EnvString("RABBITMQ_USER", "guest")
	AMQP_PASS    = common.EnvString("RABBITMQ_PASS", "guest")
	AMQP_HOST    = common.EnvString("RABBITMQ_HOST", "localhost")
	AMQP_PORT    = common.EnvString("RABBITMQ_PORT", "5672")
	STRIPE_KEY   = common.EnvString("STRIPE_KEY", "sk_test_51Pkl3SFmQPISRdlo8SIHtHfspZaAW31DKkhAVpbuCbWhNpo7SPp6rIih11iAQK4UacX0Q9zUAdBoMtbLUGpCBBBB00cjPn5oUy")
)

func main() {
	ctx := context.Background()
	instanceId := discovery.GenerateInstanceId(SERVICE_NAME)

	registry, err := consul.NewRegistry(CONSUL_ADDR)
	if err != nil {
		panic(err)
	}
	defer registry.Deregister(ctx, instanceId, SERVICE_NAME)

	// stripe setup
	fmt.Printf("STRIPE KEY = %s", STRIPE_KEY)
	stripe.Key = "sk_test_51Pkl3SFmQPISRdlo8SIHtHfspZaAW31DKkhAVpbuCbWhNpo7SPp6rIih11iAQK4UacX0Q9zUAdBoMtbLUGpCBBBB00cjPn5oUy"
	stripeProcessor := stripeprocessor.NewProcessor()
	svc := NewService(stripeProcessor)

	//consul reg-dereg setup
	//paymentProcessor := stripe.NewProcessor()

	if err := registry.Register(ctx, instanceId, SERVICE_NAME, GRPC_ADDR); err != nil {
		panic(err)
	}

	// messaging broker setup
	ch, fnclose := broker.Connect(AMQP_USER, AMQP_USER, AMQP_HOST, AMQP_PORT)
	defer func() {
		ch.Close()
		fnclose()
	}()

	amqpConsumer := NewConsumer(svc, ch)

	go amqpConsumer.Listen()

	grpcServer := grpc.NewServer()
	ln, err := net.Listen("tcp", GRPC_ADDR)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ln.Close()

	//NewGRPCHandler(grpcServer, ch)

	log.Printf("Payments GRPC Server listening on %s", GRPC_ADDR)

	done := make(chan bool)
	defer func() {
		close(done)
	}()

	go healthCheck(done, registry, instanceId)

	if err := grpcServer.Serve(ln); err != nil {
		log.Fatal(err.Error())
	}
}

func healthCheck(done chan bool, registry *consul.Registry, instanceId string) {
	for {
		select {
		case <-done:
			return
		default:
			if err := registry.HealthCheck(instanceId, SERVICE_NAME); err != nil {
				log.Fatalf("Failed to health check instance %s of service %s", instanceId, SERVICE_NAME)
				close(done)
			}
		}
		time.Sleep(time.Second * 2) // punctuate the calls to HealthCheck
	}
}
