package main

import (
	"log"
	"net"

	common "github.com/swarajroy/oms-common"
	"google.golang.org/grpc"
)

var (
	ORDER_SCV_GRPC_ADDR = common.EnvString("ORDER_SVC_GRPC_ADDR", "localhost:2000")
)

func main() {
	grpcServer := grpc.NewServer()
	ln, err := net.Listen("tcp", ORDER_SCV_GRPC_ADDR)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ln.Close()
	store := NewStore()
	svc := NewService(store)
	NewGRPCHandler(grpcServer, svc)

	log.Printf("Orders GRPC Server listening on %s", ORDER_SCV_GRPC_ADDR)

	if err := grpcServer.Serve(ln); err != nil {
		log.Fatal(err.Error())
	}

}
