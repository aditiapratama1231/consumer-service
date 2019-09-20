package service

import (
	"magento-consumer-service/domain"
)

type OrderService interface {
	UpdateStatusOrder(*domain.Consume) error
}
