package service

import (
	"magento-consumer-service/domain"
)

type ProductService interface {
	CreateProduct(*domain.Consume) error
	UpdateProduct(*domain.Consume) error
	DeleteProduct(*domain.Consume) error
}
