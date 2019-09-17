package repository

import (
	"magento-consumer-service/database/models"
	"magento-consumer-service/domain"
)

// CustomerRepository /
type CustomerRepository interface {
	SyncCustomer(domain.CustomerRecord) (interface{}, error)
	ShowCustomerID(domain.CustomerRecord) (interface{}, error)
	GetMagentoID(string, int) (models.CustomerRecord, error)
	SaveStream(string) (interface{}, error)
}
