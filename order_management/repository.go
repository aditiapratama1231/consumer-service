package repository

import (
	"magento-consumer-service/domain"
	"magento-consumer-service/database/models"
)

type OrderRepository interface {
	SyncOrder(domain.OrderRecord) (interface{}, error)
	GetMagentoID(string, int) (models.OrderRecord, error)
	SaveStream(string) (interface{}, error)
	DeleteRecord(string, int) error
}
