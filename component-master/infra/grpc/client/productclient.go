package client

import (
	"component-master/config"
	rpc "component-master/proto/product"
	"log/slog"

	"google.golang.org/grpc"
)

var (
	productGrpcClient *productClient
)

type productClient struct {
	client rpc.ProductServiceClient
}

func NewProductClient() *productClient {
	return productGrpcClient
}

func InitGrpcProductClient(conf config.GrpcConfigClient, logConf config.LogConfig) {
	doOne.Do(func() {
		conn, err := InitConnection(conf, logConf)
		if err != nil {
			slog.Error("cannot initial auth grpc client", "error", err)
			return
		}
		initPromotionClient(conn)
	})
}

func initProductClient(conn *grpc.ClientConn) {
	productGrpcClient = &productClient{
		client: rpc.NewProductServiceClient(conn),
	}
}

func (r *productClient) GetProductByCode(req *rpc.GetProductByIdRequest) (*rpc.GetProductByIdResponse, error) {
	ctx := ContextwithTimeout()
	resp, err := r.client.GetProductById(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *productClient) CreateProduct(req *rpc.CreateProductRequest) (*rpc.CreateProductResponse, error) {
	ctx := ContextwithTimeout()
	resp, err := r.client.CreateProduct(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
