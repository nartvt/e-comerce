package repository

import (
	"promotion-service/app/domain/model"
	"promotion-service/app/infra/db"

	"gorm.io/gorm"
)

type promotionRepository struct {
	dbInstance *gorm.DB
}

type IPromotionRepository interface {
	GetPromotionById(id int) (*model.Promotion, error)
	GetPromotionsPagination(isActive bool, limit int, offset int) ([]*model.Promotion, error)
	UpdatePromotion(promotion *model.Promotion) error
	CreatePromotion(promotion *model.Promotion) error
}

var PromotionRepository IPromotionRepository

func init() {
	PromotionRepository = &promotionRepository{dbInstance: db.GetDB()}
}

func (p promotionRepository) UpdatePromotion(promotion *model.Promotion) error {
	return p.dbInstance.Save(promotion).Error
}

func (p promotionRepository) CreatePromotion(promotion *model.Promotion) error {
	return p.dbInstance.Create(promotion).Error
}

func (p promotionRepository) GetPromotionById(id int) (*model.Promotion, error) {
	var resp *model.Promotion
	err := p.dbInstance.Model(&model.Promotion{}).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Find(resp).
		Error
	return resp, err
}

func (p promotionRepository) GetPromotionsPagination(isActive bool, limit int, offset int) ([]*model.Promotion, error) {
	var resp []*model.Promotion
	err := p.dbInstance.Model(&model.Promotion{}).Where("deleted_at IS NULL").
		Where("active = ?", isActive).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&resp).
		Error
	return resp, err
}
