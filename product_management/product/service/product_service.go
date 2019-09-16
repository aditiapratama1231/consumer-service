package service

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"

	"magento-consumer-service/config"
	"magento-consumer-service/domain"
	product "magento-consumer-service/product_management"
)

type productService struct {
	DB         *gorm.DB
	Repository product.ProductRepository
	Request    config.Request
}

func NewProductService(db *gorm.DB, repository product.ProductRepository, request config.Request) ProductService {
	return &productService{
		DB:         db,
		Repository: repository,
		Request:    request,
	}
}

func (product *productService) CreateProduct(consume *domain.Consume) error {
	var (
		dataResponse interface{}
	)

	payload := consume.Data.Body.Payload
	reqBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error encoding product payload " + err.Error())
	}

	req, err := product.Request.Post("/products", reqBody)
	if err != nil {
		log.Println("Error SetUp API call : ", err.Error())
	}

	req.ToJSON(&dataResponse)
	log.Println(dataResponse)

	_, err = product.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		return err
	}

	dashboardID, err := strconv.Atoi(consume.Data.Head.Dashboard)
	_, err = product.Repository.SyncProduct(domain.ProductRecord{
		Type:           "product",
		MagentoID:      1,
		DashboardID:    dashboardID,
		SequenceNumber: *consume.SequenceNumber,
	})

	return nil
}

func (product *productService) UpdateProduct(consume *domain.Consume) error {
	return nil
}

func (product *productService) DeleteProduct(consume *domain.Consume) error {
	return nil
}
