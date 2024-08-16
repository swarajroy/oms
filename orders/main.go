package main

import (
	"context"
	"log"
	"net"
	"time"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/swarajroy/oms-common"
	"github.com/swarajroy/oms-common/broker"
	"github.com/swarajroy/oms-common/discovery"
	"github.com/swarajroy/oms-common/discovery/consul"
	"google.golang.org/grpc"
)

var (
	SERVICE_NAME = "orders"
	GRPC_ADDR    = common.EnvString("GRPC_ADDR", "localhost:2000")
	CONSUL_ADDR  = common.EnvString("CONSUL_ADDR", "localhost:8500")
	AMQP_USER    = common.EnvString("RABBITMQ_USER", "guest")
	AMQP_PASS    = common.EnvString("RABBITMQ_PASS", "guest")
	AMQP_HOST    = common.EnvString("RABBITMQ_HOST", "localhost")
	AMQP_PORT    = common.EnvString("RABBITMQ_PORT", "5672")
)

func main() {
	ctx := context.Background()
	instanceId := discovery.GenerateInstanceId(SERVICE_NAME)

	registry, err := consul.NewRegistry(CONSUL_ADDR)
	if err != nil {
		panic(err)
	}
	defer registry.Deregister(ctx, instanceId, SERVICE_NAME)

	if err := registry.Register(ctx, instanceId, SERVICE_NAME, GRPC_ADDR); err != nil {
		panic(err)
	}

	ch, fnclose := broker.Connect(AMQP_USER, AMQP_USER, AMQP_HOST, AMQP_PORT)
	defer func() {
		ch.Close()
		fnclose()
	}()

	grpcServer := grpc.NewServer()
	ln, err := net.Listen("tcp", GRPC_ADDR)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ln.Close()
	store := NewStore()
	svc := NewService(store)

	NewGRPCHandler(grpcServer, ch, svc)

	log.Printf("Orders GRPC Server listening on %s", GRPC_ADDR)

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
