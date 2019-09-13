package service

import (
	"encoding/json"
	"errors"
	"log"
	"magento-consumer-service/config"
	"magento-consumer-service/database/models"
	"magento-consumer-service/domain"
	"os"

	"net/http"
)

// BrandService interface contain all service or function for brand
type BrandService interface {
	CreateBrand() error
	UpdateBrand() error
	DeleteBrand() error
}

// CreateBrand /
func (brand *Service) CreateBrand() error {
	var (
		kinesis      models.KinesisSequenceNumber
		record       models.ProductRecord
		client       = &http.Client{}
		dataResponse interface{}
		brandData    domain.Brand
	)

	rqst := config.NewRequest(os.Getenv("MAGENTO_BASE_URL"))

	// convert brand payload
	if brand.Consume.Data.Body.Payload["brand"] != nil {
		brandData = convertBrand(brand.Consume.Data.Body.Payload["brand"])
	} else {
		return errors.New("wrong payload")
	}

	// POST data
	reqBody, err := json.Marshal(brand.Consume.Data)
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
	kinesis.SequenceNumber = *brand.Consume.SequenceNumber
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
func (brand *Service) UpdateBrand() error {
	log.Println("update")
	log.Println(*brand.Consume.SequenceNumber)
	return nil
}

// DeleteBrand /
func (brand *Service) DeleteBrand() error {
	log.Println("delete")
	log.Println(*brand.Consume.SequenceNumber)
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
