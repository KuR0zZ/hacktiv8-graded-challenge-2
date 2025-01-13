package main

import (
	"graded-challenge-2-server/config"
	"graded-challenge-2-server/helper"
	"graded-challenge-2-server/middleware"
	"graded-challenge-2-server/pb"
	"graded-challenge-2-server/repository"
	"graded-challenge-2-server/service"
	"log"
	"net"
	"os"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, skipping....")
	}

	db := config.InitDB()

	serverRepository := repository.NewServerRepositoryImpl(db)
	serverService := service.NewServerService(serverRepository)

	cronJob := helper.NewCronJob(*serverService)
	cronJob.UpdateBookStatus()

	port := os.Getenv("PORT")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc_auth.UnaryServerInterceptor(middleware.JWTAuth),
		),
	)

	pb.RegisterServerServiceServer(grpcServer, serverService)

	log.Println("Server is running on port:", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
