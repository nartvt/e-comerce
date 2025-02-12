package grpc

import (
	"context"
	"fmt"
	"log"
	"net"
	"promotion-service/config"
	"promotion-service/infra/db"
	"promotion-service/internal/protobuf"
	"promotion-service/internal/repository"
	"promotion-service/internal/services"
	"time"

	"google.golang.org/grpc"
)

func InitGrpcServer() {
	conf := config.Config
	address := fmt.Sprintf("%s:%d", conf.Server.Host, conf.Server.Port)
	listen, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	promotionRepo := repository.NewPromotionRepository(db.GetDB())
	promotionService := services.NewPromotionService(promotionRepo)

	ops := []grpc.ServerOption{
		grpc.ConnectionTimeout(time.Duration(conf.Server.TimeOut) * time.Second), // set connection timeout
		grpc.UnaryInterceptor(interceptor),                                       // set unary interceptor
	}

	grpcServer := grpc.NewServer(ops...)
	protobuf.RegisterPromotionServiceServer(grpcServer, promotionService)
	log.Printf("Listening on %v", listen.Addr())
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func interceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {

	log.Printf("received request: %v", req)
	resp, err := handler(ctx, req)
	log.Printf("sent response: %v", resp)
	return resp, err
}
