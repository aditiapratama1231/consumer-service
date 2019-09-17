package service

import (
	"encoding/json"
	"errors"
	clog "log"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

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
		dataResponse interface{}
		brandData    domain.Brand
	)

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
		log.Error("Error Encoding brand payload : " + err.Error())
		return err
	}

	req, err := brand.Request.Post("/products", reqBody)
	if err != nil {
		clog.Printf("%+v", req) // force to show request details in command line.
		log.Error("Error SetUp API call : " + err.Error())
		return err
	}

	req.ToJSON(&dataResponse)

	// if POST success, safe data to db
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
	log.Error("update")
	// log.Error(*consume.SequenceNumber)
	return nil
}

// DeleteBrand /
func (brand *brandService) DeleteBrand(domain *domain.Consume) error {
	log.Error("delete")
	// log.Error(*consume.SequenceNumber)
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
