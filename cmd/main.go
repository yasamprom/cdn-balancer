package main

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/spf13/viper"
	"github.com/yasamprom/cdn-balancer/internal/usecases"
	"google.golang.org/grpc/reflection"

	app "github.com/yasamprom/cdn-balancer/internal/cdn-balancer"
	pb "github.com/yasamprom/cdn-balancer/internal/pb/api"
)

func main() {
	// Check if config is valid
	parseConfig()

	// Create service
	uc := usecases.New()
	balancer := app.NewBalancerService(&app.Config{
		Usecases: uc,
	})

	// Run service
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", 8080))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	reflection.Register(grpcServer)
	pb.RegisterBalancerServer(grpcServer, balancer)
	grpcServer.Serve(lis)
}

func parseConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("failed to parse config file: %w", err))
	}
}
