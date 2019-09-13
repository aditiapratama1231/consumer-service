package service

import (
	"magento-consumer-service/domain"
)

type BrandService interface {
	CreateBrand(*domain.Consume) error
	UpdateBrand(*domain.Consume) error
	DeleteBrand(*domain.Consume) error
}
