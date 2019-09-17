package service

import (
	"magento-consumer-service/domain"
)

//CustomerService /
type CustomerService interface {
	CreateCustomer(*domain.Consume) error
	UpdateCustomer(*domain.Consume) error
	DeleteCustomer(*domain.Consume) error
}
