package entity

import (
	"context"
	"fmt"
	"promotion-service/app/domain/model"
	"promotion-service/app/infra/db"
	"promotion-service/app/protogen/message"
	"promotion-service/app/protogen/rpc"
	"promotion-service/app/uerror"
	"promotion-service/app/usercases/promotion/repository"
	"time"

	"gorm.io/gorm"
)

type prodmotionEntity struct {
	rpc.UnimplementedPromotionServiceServer
}

var PromotionEntity rpc.PromotionServiceServer

func init() {
	PromotionEntity = &prodmotionEntity{
		UnimplementedPromotionServiceServer: rpc.UnimplementedPromotionServiceServer{},
	}
}

func (p prodmotionEntity) GetPromotionById(ctx context.Context, request *message.GetPromotionByIdRequest) (*message.GetPromotionByIdResponse, error) {
	if request == nil || request.Id <= 0 {
		return &message.GetPromotionByIdResponse{}, nil
	}
	promotion, err := repository.PromotionRepository.GetPromotionById(int(request.Id), true)
	if err != nil && err == gorm.ErrRecordNotFound {
		return &message.GetPromotionByIdResponse{}, nil
	}
	if err != nil {
		return &message.GetPromotionByIdResponse{}, uerror.InternalError(err, err.Error())
	}
	if promotion == nil {
		return &message.GetPromotionByIdResponse{}, uerror.BadRequestError(fmt.Errorf("promotion not found %d", request.Id), "promotion not found")
	}
	return &message.GetPromotionByIdResponse{Promotion: toPromotionMessage(promotion)}, nil
}

func (p prodmotionEntity) CreatePromotion(ctx context.Context, req *message.CreatePromotionRequest) (*message.CreatePromotionResponse, error) {
	if req == nil {
		return &message.CreatePromotionResponse{}, nil
	}

	modelToBeSaved := &model.Promotion{
		Code:                   req.Code,
		PromotionType:          req.PromotionType,
		Value:                  int(req.Value),
		CreatedBy:              int(req.CreatedBy),
		CreatedAt:              time.Unix(req.CreatedAt, 0).Local(),
		UpdatedAt:              time.Unix(req.UpdatedAt, 0).Local(),
		DeletedAt:              time.Unix(req.DeletedAt, 0).Local(),
		ActiveFrom:             time.Unix(req.ActiveFrom, 0).Local(),
		ActiveTo:               time.Unix(req.ActiveTo, 0).Local(),
		DailyActiveFrom:        int(req.DailyActiveFrom),
		DailyActiveTo:          int(req.DailyActiveTo),
		MaxActiveTime:          int(req.MaxActiveTime),
		MaxDailyActiveTime:     int(req.MaxDailyActiveTime),
		PerUserActiveTime:      int(req.PerUserActiveTime),
		PerUserDailyActiveTime: int(req.PerUserDailyActiveTime),
		Active:                 req.Active,
	}
	err := repository.PromotionRepository.CreatePromotionTx(modelToBeSaved, db.DB())
	if err != nil {
		return &message.CreatePromotionResponse{}, err
	}

	return &message.CreatePromotionResponse{Promotion: &message.Promotion{Id: int64(modelToBeSaved.Id)}}, nil
}

func (p prodmotionEntity) UpdatePromotion(ctx context.Context, req *message.UpdatePromotionRequest) (*message.UpdatePromotionResponse, error) {
	if req == nil || req.Id <= 0 {
		return &message.UpdatePromotionResponse{}, nil
	}

	err := repository.PromotionRepository.UpdatePromotionTx(&model.Promotion{
		Code:                   req.Code,
		PromotionType:          req.PromotionType,
		Value:                  int(req.Value),
		CreatedBy:              int(req.CreatedBy),
		CreatedAt:              time.Unix(req.CreatedAt, 0).Local(),
		UpdatedAt:              time.Unix(req.UpdatedAt, 0).Local(),
		DeletedAt:              time.Unix(req.DeletedAt, 0).Local(),
		ActiveFrom:             time.Unix(req.ActiveFrom, 0).Local(),
		ActiveTo:               time.Unix(req.ActiveTo, 0).Local(),
		DailyActiveFrom:        int(req.DailyActiveFrom),
		DailyActiveTo:          int(req.DailyActiveTo),
		MaxActiveTime:          int(req.MaxActiveTime),
		MaxDailyActiveTime:     int(req.MaxDailyActiveTime),
		PerUserActiveTime:      int(req.PerUserActiveTime),
		PerUserDailyActiveTime: int(req.PerUserDailyActiveTime),
		Active:                 req.Active,
	}, db.DB())
	if err != nil {
		return &message.UpdatePromotionResponse{}, err
	}

	return &message.UpdatePromotionResponse{Promotion: &message.Promotion{Id: req.Id}}, nil
}

func toPromotionMessage(p *model.Promotion) *message.Promotion {
	return &message.Promotion{
		Id:                     int64(p.Id),
		Code:                   p.Code,
		PromotionType:          p.PromotionType,
		Value:                  int64(p.Value),
		CreatedBy:              int64(p.CreatedBy),
		CreatedAt:              p.CreatedAt.Unix(),
		UpdatedAt:              p.UpdatedAt.Unix(),
		DeletedAt:              p.DeletedAt.Unix(),
		ActiveFrom:             p.ActiveFrom.Unix(),
		ActiveTo:               p.ActiveTo.Unix(),
		DailyActiveFrom:        int64(p.DailyActiveFrom),
		DailyActiveTo:          int64(p.DailyActiveTo),
		MaxActiveTime:          int32(p.MaxActiveTime),
		MaxDailyActiveTime:     int32(p.MaxDailyActiveTime),
		PerUserActiveTime:      int32(p.PerUserActiveTime),
		PerUserDailyActiveTime: int32(p.PerUserDailyActiveTime),
		Active:                 p.Active,
	}
}
