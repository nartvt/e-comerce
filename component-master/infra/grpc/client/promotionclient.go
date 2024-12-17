package client

import (
	"component-master/config"
	rpc "component-master/proto/promotion"
	"log/slog"

	"google.golang.org/grpc"
)

var (
	promotionGrpcClient *promotionClient
)

type promotionClient struct {
	client rpc.PromotionServiceClient
}

func NewPromotionClient() *promotionClient {
	return promotionGrpcClient
}

func InitGrpcPromotionClient(conf config.GrpcConfigClient, logConf config.LogConfig) {
	doOne.Do(func() {
		conn, err := InitConnection(conf, logConf)
		if err != nil {
			slog.Error("cannot initial auth grpc client", "error", err)
			return
		}
		initPromotionClient(conn)
	})
}

func initPromotionClient(conn *grpc.ClientConn) {
	promotionGrpcClient = &promotionClient{
		client: rpc.NewPromotionServiceClient(conn),
	}
}

func (r *promotionClient) GetPromotionByCode(req *rpc.GetProductionByCodeRequest) (*rpc.GetPromotionByCodeResponse, error) {
	ctx := ContextwithTimeout()
	resp, err := r.client.GetPromotionByCode(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *promotionClient) CreatePromotion(req *rpc.PromotionCreate) (*rpc.CreatePromotionResponse, error) {
	ctx := ContextwithTimeout()
	resp, err := r.client.CreatePromotion(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
