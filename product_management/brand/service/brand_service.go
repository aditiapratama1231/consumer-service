package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"

	"magento-consumer-service/config"
	"magento-consumer-service/database/models"
	"magento-consumer-service/domain"
	product "magento-consumer-service/product_management"
)

type brandService struct {
	DB         *gorm.DB
	Repository product.ProductRepository
	Request    config.Request
}

func NewBrandService(db *gorm.DB, repository product.ProductRepository, request config.Request) BrandService {
	return &brandService{
		DB:         db,
		Repository: repository,
		Request:    request,
	}
}

// CreateBrand /
func (brand *brandService) CreateBrand(consume *domain.Consume) error {
	var (
		kinesis      models.KinesisSequenceNumber
		record       models.ProductRecord
		client       = &http.Client{}
		dataResponse interface{}
		brandData    domain.Brand
	)

	rqst := config.NewRequest(os.Getenv("MAGENTO_BASE_URL"))

	// convert brand payload
	if consume.Data.Body.Payload["brand"] != nil {
		brandData = convertBrand(consume.Data.Body.Payload["brand"])
	} else {
		return errors.New("wrong payload")
	}

	// POST data
	reqBody, err := json.Marshal(consume.Data)
	if err != nil {
		return err
	}

	request, err := rqst.Post("/products", reqBody)
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&dataResponse)
	if err != nil {
		return err
	}

	// if POST success, safe data to db
	kinesis.SequenceNumber = *consume.SequenceNumber
	err = brand.DB.Create(&kinesis).Error
	if err != nil {
		return err
	}

	// save to record
	record.MagentoID = 1 // change this code next
	record.DashboardID = brandData.ID
	err = brand.DB.Create(&record).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateBrand /
func (brand *brandService) UpdateBrand(domain *domain.Consume) error {
	log.Println("update")
	// log.Println(*consume.SequenceNumber)
	return nil
}

// DeleteBrand /
func (brand *brandService) DeleteBrand(domain *domain.Consume) error {
	log.Println("delete")
	// log.Println(*consume.SequenceNumber)
	return nil
}

func convertBrand(data interface{}) domain.Brand {
	m := data.(map[string]interface{})
	brand := domain.Brand{}
	if name, ok := m["name"].(string); ok {
		brand.Name = name
	}

	if id, ok := m["id"].(float64); ok {
		brand.ID = int(id)
	}
	return brand
}
