package server

import (
	"component-master/config"
	"component-master/proto/product"
)

func InitGrpcProductServer(conf config.ServerInfo, productServer product.ProductServiceServer) *GrpcServer {
	grpcServer := &GrpcServer{
		conf: &conf,
	}
	grpcServer.InitGrpcServer()
	product.RegisterProductServiceServer(grpcServer.server, productServer)
	return grpcServer
}
