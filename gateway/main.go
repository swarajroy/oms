package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	common "github.com/swarajroy/oms-common"
	pb "github.com/swarajroy/oms-common/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	HTTP_ADDR          = common.EnvString("HTTP_ADDR", ":8080")
	ORDER_SERVICE_ADDR = "localhost:2000"
)

func main() {

	conn, err := grpc.NewClient(ORDER_SERVICE_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to dial a connection with GRPC")
	}
	defer conn.Close()

	log.Printf("Dialing connection to orders service at %s", ORDER_SERVICE_ADDR)

	c := pb.NewOrderServiceClient(conn)

	mux := http.NewServeMux()
	handler := NewHandler(c)
	handler.RegisterRoutes(mux)

	log.Printf("Starting the http server on port %s", HTTP_ADDR)

	if err := http.ListenAndServe(HTTP_ADDR, mux); err != nil {
		log.Fatalf("Failed to start http server on port %s", HTTP_ADDR)
	}
}
