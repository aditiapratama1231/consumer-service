package product

import (
	"github.com/jinzhu/gorm"

	"magento-consumer-service/database/models"
)

type ProductRepository interface {
	SyncProduct(int64, int64) (models.ProductRecord, error)
	ShowProductID(int64, int64) (models.ProductRecord, error)
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) ProductRepository {
	return &productRepository{
		DB: DB,
	}
}

func (p *productRepository) SyncProduct(dashboardID int64, magentoID int64) (models.ProductRecord, error) {
	var product models.ProductRecord

	return product, nil
}

func (p *productRepository) ShowProductID(dashboardID int64, magentoID int64) (models.ProductRecord, error) {
	var product models.ProductRecord

	return product, nil
}
