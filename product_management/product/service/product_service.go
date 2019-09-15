package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

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
		productData  domain.Product
	)
	fmt.Println("hereee")
	payload := consume.Data.Body.Payload["product"]
	if payload != nil {
		productData = convertProduct(payload)
	} else {
		return errors.New("wrong payload")
	}

	reqBody, err := json.Marshal(consume.Data)
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

	_, err = product.Repository.SyncProduct(domain.ProductRecord{
		Type:        "product",
		MagentoID:   1,
		DashboardID: productData.ProductSupplierID,
	})

	return nil
}

func (product *productService) UpdateProduct(consume *domain.Consume) error {
	return nil
}

func (product *productService) DeleteProduct(consume *domain.Consume) error {
	return nil
}

func convertProduct(data interface{}) domain.Product {
	m := data.(map[string]interface{})
	product := domain.Product{}
	if name, ok := m["name"].(string); ok {
		product.Name = name
	}

	if id, ok := m["id"].(float64); ok {
		product.ProductSupplierID = int(id)
	}
	return product
}
