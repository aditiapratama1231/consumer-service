package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"

	"magento-consumer-service/config"
	"magento-consumer-service/database/models"
	"magento-consumer-service/domain"
	"magento-consumer-service/product"
	"magento-consumer-service/product/category"
)

type categoryService struct {
	DB         *gorm.DB
	Consume    *domain.Consume
	Repository *product.ProductRepository
}

func NewCategoryService(DB *gorm.DB, domain *domain.Consume) category.CategoryService {
	return &categoryService{
		DB:      DB,
		Consume: domain,
	}
}

func (c *categoryService) CreateCategory() error {
	var (
		kinesis      models.KinesisSequenceNumber
		record       models.ProductRecord
		client       = &http.Client{}
		dataResponse interface{}
		categoryData domain.Category
	)

	request := config.NewRequest(os.Getenv("MAGENTO_BASE_URL"))

	// convert category payload
	if c.Consume.Data.Body.Payload["category"] != nil {
		categoryData = convertCategory(c.Consume.Data.Body.Payload["category"])
	} else {
		return errors.New("wrong payload")
	}

	reqBody, err := json.Marshal(c.Consume.Data)
	if err != nil {
		return err
	}

	// POST Data
	req, err := request.Post("/products", reqBody)
	if err != nil {
		return err
	}

	response, err := client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&dataResponse)
	if err != nil {
		return err
	}

	// if POST success, safe data to db
	kinesis.SequenceNumber = *c.Consume.SequenceNumber
	err = c.DB.Create(&kinesis).Error
	if err != nil {
		return err
	}

	// save to record
	record.MagentoID = 1 // change this code next
	record.DashboardID = categoryData.ID
	err = c.DB.Create(&record).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *categoryService) UpdateCategory() error {
	return nil
}

func (c *categoryService) DeleteCategory() error {
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
