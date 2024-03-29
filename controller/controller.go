package controller

import (
	"magento-consumer-service/database/models"
	"magento-consumer-service/domain"

	log "github.com/sirupsen/logrus"

	_customerService "magento-consumer-service/customer_management/customer/service"
	_orderService "magento-consumer-service/order_management/order/service"
	_brandService "magento-consumer-service/product_management/brand/service"
	_categoryService "magento-consumer-service/product_management/category/service"
	_productService "magento-consumer-service/product_management/product/service"

	"github.com/jinzhu/gorm"
)

// Controller structs /
type controller struct {
	DB              *gorm.DB
	CategoryService _categoryService.CategoryService
	BrandService    _brandService.BrandService
	ProductService  _productService.ProductService
	CustomerService _customerService.CustomerService
	OrderService    _orderService.OrderService
}

//Controller /
type Controller interface {
	MainController(*domain.Consume)
	ProductManagement(*domain.Consume)
	OrderManagement(*domain.Consume)
	CustomerManagement(*domain.Consume)
}

func NewController(db *gorm.DB,
	categorySrv _categoryService.CategoryService,
	brandSrv _brandService.BrandService,
	productSrv _productService.ProductService,
	customerSrv _customerService.CustomerService,
	orderSrv _orderService.OrderService,
) Controller {
	return &controller{
		DB:              db,
		CategoryService: categorySrv,
		BrandService:    brandSrv,
		ProductService:  productSrv,
		CustomerService: customerSrv,
		OrderService:    orderSrv,
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
		case "merchant":
			c.CustomerManagement(consume)
		default:
			log.Fatal("wrong service input")
		}
	}
}
