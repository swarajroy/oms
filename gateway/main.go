package main

import (
	"context"
	"log"
	"net/http"
	"time"

	
	co_ "github.com/joho/godotenv/autoload"mmon "github.com/swarajroy/oms-common"
	"github.com/swarajroy/oms-common/discovery"
	"github.com/swarajroy/oms-common/discovery/consul"
	"github.com/swarajroy/oms-gateway/gateway"
)

var (
	SERVICE_NAME = "gateway"
	HTTP_ADDR    = common.EnvString("HTTP_ADDR", ":8080")
	CONSUL_ADDR  = common.EnvString("CONSUL_ADDR", "localhost:8500")
)

func main() {

	ctx := context.Background()
	instanceId := discovery.GenerateInstanceId(SERVICE_NAME)

	registry, err := consul.NewRegistry(CONSUL_ADDR)
	if err != nil {
		panic(err)
	}
	defer registry.Deregister(ctx, instanceId, SERVICE_NAME)

	if err := registry.Register(ctx, instanceId, SERVICE_NAME, HTTP_ADDR); err != nil {
		panic(err)
	}

	grpcGateway := gateway.NewGRPCGateway(registry)
	handler := NewHandler(grpcGateway)

	log.Printf("Starting the http server on port %s", HTTP_ADDR)

	done := make(chan bool)
	defer close(done)

	go healthCheck(done, registry, instanceId)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	if err := http.ListenAndServe(HTTP_ADDR, mux); err != nil {
		log.Fatalf("Failed to start http server on port %s", HTTP_ADDR)
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
