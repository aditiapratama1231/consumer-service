package repository

import (
	"magento-consumer-service/database/models"
	"magento-consumer-service/domain"
	product "magento-consumer-service/product_management"

	"github.com/jinzhu/gorm"
)

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) product.ProductRepository {
	return &productRepository{
		DB: DB,
	}
}

func (p *productRepository) SyncProduct(data domain.ProductRecord) (interface{}, error) {
	var product models.ProductRecord

	product.Type = data.Type
	product.DashboardID = data.DashboardID
	product.MagentoID = data.MagentoID
	product.SequenceNumber = data.SequenceNumber

	tx := p.DB.Begin()

	err := tx.Create(&product).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return product, nil
}

func (p *productRepository) ShowProductID(data domain.ProductRecord) (interface{}, error) {
	var product models.ProductRecord

	p.DB.Where("type = ?", data.Type).Where("magento_id = ?", data.MagentoID).Where("dashboard_id = ?", data.DashboardID).Find(&product)

	return product, nil
}

func (p *productRepository) SaveStream(sequenceNumber string) (interface{}, error) {
	var kinesis models.KinesisSequenceNumber

	kinesis.SequenceNumber = sequenceNumber

	tx := p.DB.Begin()

	err := tx.Create(&kinesis).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return kinesis, nil
}
