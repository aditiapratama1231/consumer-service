package repository

import (
	"magento-consumer-service/domain"
)

type ProductRepository interface {
	SyncProduct(domain.ProductRecord) (interface{}, error)
	ShowProductID(domain.ProductRecord) (interface{}, error)
	SaveStream(string) (interface{}, error)
}
