package service

import (
	"magento-consume-service/domain"

	"github.com/jinzhu/gorm"
)

// Service structs /
type Service struct {
	DB      *gorm.DB
	Consume *domain.Consume
}
