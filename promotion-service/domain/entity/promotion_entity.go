package entity

import (
	"context"
	"errors"
	"promotion-service/domain/model"
	"promotion-service/protogen/message"
	"promotion-service/protogen/rpc"
	"promotion-service/repository"
	"promotion-service/uerror"
	"time"

	"gorm.io/gorm"
)

type promotionEntity struct {
	rpc.UnimplementedPromotionServiceServer
}

var PromotionEntity rpc.PromotionServiceServer

func init() {
	PromotionEntity = &promotionEntity{
		UnimplementedPromotionServiceServer: rpc.UnimplementedPromotionServiceServer{},
	}
}

func (p promotionEntity) GetPromotionById(ctx context.Context, request *message.GetPromotionByIdRequest) (*message.GetPromotionByIdResponse, error) {
	if request == nil || request.Id <= 0 {
		err := errors.New("request is empty")
		return nil, uerror.BadRequestError(err)
	}

	promotion, err := repository.PromotionRepository.GetPromotionById(int(request.Id))
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, uerror.NotFoundError(err)
	}
	if err != nil {
		return nil, uerror.InternalError(err, err.Error())
	}

	return &message.GetPromotionByIdResponse{Promotion: toPromotionMessage(promotion)}, nil
}

func (p promotionEntity) CreatePromotion(ctx context.Context, req *message.CreatePromotionRequest) (*message.CreatePromotionResponse, error) {
	if req == nil {
		err := errors.New("request is empty")
		return nil, uerror.BadRequestError(err)
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
	err := repository.PromotionRepository.CreatePromotion(modelToBeSaved)
	if err != nil {
		return nil, uerror.InternalError(err)
	}

	return &message.CreatePromotionResponse{}, nil
}

func (p promotionEntity) UpdatePromotion(ctx context.Context, req *message.UpdatePromotionRequest) (*message.UpdatePromotionResponse, error) {
	if req == nil || req.Id <= 0 {
		err := errors.New("request is empty")
		return nil, uerror.BadRequestError(err)
	}

	err := repository.PromotionRepository.UpdatePromotion(&model.Promotion{
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
	})
	if err != nil {
		return nil, uerror.InternalError(err)
	}

	return &message.UpdatePromotionResponse{}, nil
}

func (p promotionEntity) GetPromotionsPagination(ctx context.Context, req *message.GetPromotionsPaginationRequest) (*message.GetPromotionsPaginationResponse, error) {
	if req == nil {
		err := errors.New("request is empty")
		return nil, uerror.BadRequestError(err)
	}

	if req.Offset < 0 || req.Limit <= 0 {
		err := errors.New("request param is invalid")
		return nil, uerror.BadRequestError(err)
	}

	promotions, err := repository.PromotionRepository.GetPromotionsPagination(req.IsActive, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, uerror.InternalError(err)
	}

	return &message.GetPromotionsPaginationResponse{
		Promotions: toPromotionMessages(promotions),
	}, nil
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

func toPromotionMessages(models []*model.Promotion) []*message.Promotion {
	res := make([]*message.Promotion, len(models))
	for i, m := range models {
		res[i] = toPromotionMessage(m)
	}

	return res
}
