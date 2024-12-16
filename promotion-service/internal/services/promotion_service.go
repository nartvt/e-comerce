package services

import (
	"context"
	"errors"
	"promotion-service/internal/models"
	"promotion-service/internal/protobuf"
	"promotion-service/internal/repository"
	"promotion-service/internal/uerror"
	"time"

	"gorm.io/gorm"
)

type promotionService struct {
	promotionRepo repository.IPromotionRepository
	protobuf.UnimplementedPromotionServiceServer
}

func NewPromotionService(promotionRepo repository.IPromotionRepository) protobuf.PromotionServiceServer {
	return promotionService{promotionRepo: promotionRepo}
}

func (p promotionService) CreatePromotion(ctx context.Context, req *protobuf.CreatePromotionRequest) (*protobuf.CreatePromotionResponse, error) {
	if req == nil {
		err := errors.New("request is empty")
		return nil, uerror.BadRequestError(err)
	}

	modelToBeSaved := &models.Promotion{
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
	err := p.promotionRepo.CreatePromotion(modelToBeSaved)
	if err != nil {
		return nil, uerror.InternalError(err)
	}

	return &protobuf.CreatePromotionResponse{}, nil
}

func (p promotionService) GetPromotionById(ctx context.Context, request *protobuf.GetPromotionByIdRequest) (*protobuf.GetPromotionByIdResponse, error) {
	if request == nil || request.Id <= 0 {
		err := errors.New("request is empty")
		return nil, uerror.BadRequestError(err)
	}

	promotion, err := p.promotionRepo.GetPromotionById(int(request.Id))
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil, uerror.NotFoundError(err)
	}
	if err != nil {
		return nil, uerror.InternalError(err, err.Error())
	}

	return &protobuf.GetPromotionByIdResponse{Promotion: toPromotionMessage(promotion)}, nil
}

func (p promotionService) GetPromotionsPagination(ctx context.Context, req *protobuf.GetPromotionsPaginationRequest) (*protobuf.GetPromotionsPaginationResponse, error) {
	if req == nil {
		err := errors.New("request is empty")
		return nil, uerror.BadRequestError(err)
	}

	if req.Offset < 0 || req.Limit <= 0 {
		err := errors.New("request param is invalid")
		return nil, uerror.BadRequestError(err)
	}

	promotions, err := p.promotionRepo.GetPromotionsPagination(req.IsActive, int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, uerror.InternalError(err)
	}

	return &protobuf.GetPromotionsPaginationResponse{
		Promotions: toPromotionMessages(promotions),
	}, nil
}

func (p promotionService) UpdatePromotion(ctx context.Context, req *protobuf.UpdatePromotionRequest) (*protobuf.UpdatePromotionResponse, error) {
	if req == nil || req.Id <= 0 {
		err := errors.New("request is empty")
		return nil, uerror.BadRequestError(err)
	}

	err := p.promotionRepo.UpdatePromotion(&models.Promotion{
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

	return &protobuf.UpdatePromotionResponse{}, nil
}

func toPromotionMessage(p *models.Promotion) *protobuf.Promotion {
	return &protobuf.Promotion{
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

func toPromotionMessages(models []*models.Promotion) []*protobuf.Promotion {
	res := make([]*protobuf.Promotion, len(models))
	for i, m := range models {
		res[i] = toPromotionMessage(m)
	}

	return res
}
