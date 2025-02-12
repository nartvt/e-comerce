package orm

import (
	"gorm.io/gorm"

	"order-service/app/domain/entities"
	"order-service/app/domain/usercases/common"
	"order-service/app/infra/db"
)

type IProduct interface {
	CreateProductTx(newProduct *entities.Product, tx *gorm.DB) error
	GetProductBySectionId(sectionId int, limit int, offset int) ([]entities.Product, int, error)
}
type product struct{}

var Product IProduct

func init() {
	Product = product{}
}
func (d product) CreateProductTx(newProduct *entities.Product, tx *gorm.DB) error {
	return tx.Save(newProduct).Error
}

func (d product) GetProductBySectionId(sectionId int, limit int, offset int) ([]entities.Product, int, error) {
	var resp []entities.Product
	total := int64(0)
	err := db.DB().Model(&entities.Product{}).
		InnerJoins("JOIN section_products ON promotion.id = section_products.product_id").
		InnerJoins("JOIN newsfeed_sections ON section_products.section_id = newsfeed_sections.id").
		Where("newsfeed_sections.active = TRUE").
		Where("promotion.active = TRUE").
		Where("type = ?", common.NewsfeedSectionTypeTopPage).
		Where("newsfeed_sections.id = ?", sectionId).
		Limit(limit).
		Offset(offset).
		Order("id DESC").
		Count(&total).
		Find(&resp).Error
	return resp, int(total), err
}
