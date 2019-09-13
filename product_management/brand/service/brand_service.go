package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"

	"magento-consumer-service/config"
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
		client       = &http.Client{}
		dataResponse interface{}
		brandData    domain.Brand
	)

	rqst := config.NewRequest(os.Getenv("MAGENTO_BASE_URL"))

	// convert brand payload
	payload := consume.Data.Body.Payload["brand"]
	if payload != nil {
		brandData = convertBrand(payload)
	} else {
		return errors.New("wrong payload")
	}

	// POST data
	reqBody, err := json.Marshal(consume.Data)
	if err != nil {
		log.Println("Error Encoding brand payload : " + err.Error())
	}

	request, err := rqst.Post("/products", reqBody)
	if err != nil {
		log.Println("Error SetUp API call : " + err.Error())
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("Error Decode Response : " + err.Error())
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&dataResponse)
	if err != nil {
		return err
	}

	// if POST success, safe data to db
	// kinesis.SequenceNumber = *consume.SequenceNumber
	// err = brand.DB.Create(&kinesis).Error
	_, err = brand.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		return err
	}

	// save to record
	_, err = brand.Repository.SyncProduct(domain.ProductRecord{
		Type:        "brand",
		MagentoID:   1,
		DashboardID: brandData.ID,
	})
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
