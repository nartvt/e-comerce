package repository

import (
	"promotion-service/internal/models"

	"gorm.io/gorm"
)

type IPromotionRepository interface {
	CreatePromotion(promotion *models.Promotion) error
	GetPromotionById(id int) (*models.Promotion, error)
	GetPromotionsPagination(isActive bool, limit int, offset int) ([]*models.Promotion, error)
	UpdatePromotion(promotion *models.Promotion) error
}

type promotionRepository struct {
	dbInstance *gorm.DB
}

func NewPromotionRepository(dbInstance *gorm.DB) IPromotionRepository {
	return promotionRepository{dbInstance: dbInstance}
}

func (p promotionRepository) CreatePromotion(promotion *models.Promotion) error {
	return p.dbInstance.Create(promotion).Error
}

func (p promotionRepository) GetPromotionById(id int) (*models.Promotion, error) {
	var resp *models.Promotion
	err := p.dbInstance.Model(&models.Promotion{}).
		Where("id = ?", id).
		Where("deleted_at IS NULL").
		Find(resp).
		Error
	return resp, err
}

func (p promotionRepository) GetPromotionsPagination(isActive bool, limit int, offset int) ([]*models.Promotion, error) {
	var resp []*models.Promotion
	err := p.dbInstance.Model(&models.Promotion{}).Where("deleted_at IS NULL").
		Where("active = ?", isActive).
		Order("id DESC").
		Offset(offset).
		Limit(limit).
		Find(&resp).
		Error
	return resp, err
}

func (p promotionRepository) UpdatePromotion(promotion *models.Promotion) error {
	return p.dbInstance.Save(promotion).Error
}
