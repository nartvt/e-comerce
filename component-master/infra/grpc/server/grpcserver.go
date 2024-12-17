package server

import (
	"component-master/config"
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	server  *grpc.Server
	address string
	conf    *config.ServerInfo
}

func (r *GrpcServer) Start() {
	if r == nil {
		slog.Error("grpc server is nil")
		return
	}
	listen, err := net.Listen("tcp", r.address)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to listen: %v", err))
		return
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := r.server.Serve(listen); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	log.Printf("Listening on %v", listen.Addr())
	<-shutdownChan // wait for the shutdown signal
	log.Println("Shutting down gRPC server...")

	r.server.GracefulStop()
	log.Println("gRPC server gracefully stopped.")
}

func (r *GrpcServer) InitGrpcServer() {
	if r == nil || r.conf == nil {
		slog.Error("grpc server is nil")
		return
	}
	r.address = fmt.Sprintf("%s:%d", r.conf.Host, r.conf.Port)
	opts := []grpc.ServerOption{
		grpc.ConnectionTimeout(time.Duration(r.conf.ConnectTimeOut) * time.Millisecond), // set connection timeout
		grpc.UnaryInterceptor(interceptor),                                              // set unary interceptor
	}
	r.server = grpc.NewServer(opts...)
	if r.conf.EnableTLS {
		// TODO: implement tls
	}
}

func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("received request: %v", req)
	log.Println(info)
	resp, err := handler(ctx, req)
	return resp, err
}
