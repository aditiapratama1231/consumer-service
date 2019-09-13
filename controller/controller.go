package controller

import (
	"log"
	"magento-consumer-service/database/models"
	"magento-consumer-service/domain"
	_brandService "magento-consumer-service/product_management/brand/service"
	_categoryService "magento-consumer-service/product_management/category/service"

	"github.com/jinzhu/gorm"
)

// Controller structs /
type controller struct {
	DB              *gorm.DB
	CategoryService _categoryService.CategoryService
	BrandService    _brandService.BrandService
}

type Controller interface {
	MainController(*domain.Consume)
	ProductManagement(*domain.Consume)
	OrderManagement(*domain.Consume)
}

func NewController(db *gorm.DB, categorySrv _categoryService.CategoryService, brandSrv _brandService.BrandService) Controller {
	return &controller{
		DB:              db,
		CategoryService: categorySrv,
		BrandService:    brandSrv,
	}
}

// MainController function, as main controller, handle incoming data record
func (c *controller) MainController(consume *domain.Consume) {
	var sequenceNumber models.KinesisSequenceNumber
	/*
		check if sequence number already exist in database.
		If true (no error), it was indicate thats data stream already consume.
		So will continue to check next data stream.
	*/
	if err := c.DB.Where("sequence_number = ?", *consume.SequenceNumber).First(&sequenceNumber).Error; err != nil {
		switch service := consume.Data.Head.Service; service {
		case "product":
			c.ProductManagement(consume)
		case "order":
			c.OrderManagement(consume)
		default:
			log.Println("wrong service input")
		}
	}
}
