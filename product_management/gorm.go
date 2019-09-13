package product

import (
	"github.com/jinzhu/gorm"

	"magento-consumer-service/database/models"
	"magento-consumer-service/domain"
)

type ProductRepository interface {
	SyncProduct(domain.ProductRecord) (interface{}, error)
	ShowProductID(domain.ProductRecord) (interface{}, error)
	SaveStream(string) (interface{}, error)
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) ProductRepository {
	return &productRepository{
		DB: DB,
	}
}

func (p *productRepository) SyncProduct(data domain.ProductRecord) (interface{}, error) {
	var product models.ProductRecord

	product.Type = data.Type
	product.DashboardID = data.DashboardID
	product.MagentoID = data.MagentoID

	err := p.DB.Create(&product).Error
	if err != nil {
		return nil, err
	}

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

	err := p.DB.Create(&kinesis).Error
	if err != nil {
		return nil, err
	}

	return kinesis, nil
}
