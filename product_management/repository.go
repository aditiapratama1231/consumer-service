package repository

import (
	"magento-consumer-service/database/models"
	"magento-consumer-service/domain"
)

type ProductRepository interface {
	SyncProduct(domain.ProductRecord) (interface{}, error)
	// UpdateProduct(domain.ProductRecord) (interface{}, error)
	ShowProductID(domain.ProductRecord) (interface{}, error)
	GetMagentoID(int) (models.ProductRecord, error)
	SaveStream(string) (interface{}, error)
}
