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

// CategoryService interface contain all service or function for brand
type CategoryService interface {
	CreateCategory() error
	UpdateCategory() error
	DeleteCategory() error
}

// CreateCategory /
func (category *Service) CreateCategory() error {
	var (
		kinesis      models.KinesisSequenceNumber
		record       models.ProductRecord
		client       = &http.Client{}
		dataResponse interface{}
		categoryData domain.Category
	)

	rqst := config.NewRequest(os.Getenv("MAGENTO_BASE_URL"))

	// convert category payload
	if category.Consume.Data.Body.Payload["category"] != nil {
		categoryData = convertCategory(category.Consume.Data.Body.Payload["category"])
	} else {
		return errors.New("wrong payload")
	}

	// POST data
	reqBody, err := json.Marshal(category.Consume.Data)
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
	kinesis.SequenceNumber = *category.Consume.SequenceNumber
	err = category.DB.Create(&kinesis).Error
	if err != nil {
		return err
	}

	// save to record
	record.MagentoID = 1 // change this code next
	record.DashboardID = categoryData.ID
	err = category.DB.Create(&record).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateCategory /
func (category *Service) UpdateCategory() error {
	log.Println("update")
	log.Println(*category.Consume.SequenceNumber)
	return nil
}

// DeleteCategory /
func (category *Service) DeleteCategory() error {
	log.Println("delete")
	log.Println(*category.Consume.SequenceNumber)
	return nil
}

func convertCategory(data interface{}) domain.Category {
	m := data.(map[string]interface{})
	category := domain.Category{}
	if name, ok := m["name"].(string); ok {
		category.Name = name
	}

	if id, ok := m["id"].(float64); ok {
		category.ID = int(id)
	}
	return category
}
