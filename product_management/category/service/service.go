package service

import (
	"magento-consumer-service/domain"
)

type CategoryService interface {
	CreateCategory(*domain.Consume) error
	UpdateCategory(*domain.Consume) error
	DeleteCategory(*domain.Consume) error
}
