package service

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/jinzhu/gorm"

	"magento-consumer-service/config"
	"magento-consumer-service/domain"
	product "magento-consumer-service/product_management"
)

type categoryService struct {
	DB         *gorm.DB
	Repository product.ProductRepository
	Request    config.Request
}

func NewCategoryService(db *gorm.DB, repository product.ProductRepository, request config.Request) CategoryService {
	return &categoryService{
		DB:         db,
		Repository: repository,
		Request:    request,
	}
}

func (c *categoryService) CreateCategory(consume *domain.Consume) error {
	var (
		client       = &http.Client{}
		dataResponse interface{}
		categoryData domain.Category
	)

	// convert category payload
	payload := consume.Data.Body.Payload["category"]
	if payload != nil {
		categoryData = convertCategory(payload)
	} else {
		return errors.New("wrong payload")
	}

	reqBody, err := json.Marshal(consume.Data)
	if err != nil {
		log.Println("Error Encoding category payload : " + err.Error())
	}

	// POST Data
	req, err := c.Request.Post("/products", reqBody)
	if err != nil {
		log.Println("Error SetUp API call : " + err.Error())
	}

	response, err := client.Do(req)
	if err != nil {
		log.Println("Error API call : " + err.Error())
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&dataResponse)
	if err != nil {
		log.Println("Error Decode Response : " + err.Error())
	}

	// if POST success, safe data to db
	_, err = c.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		return err
	}

	// save to record
	_, err = c.Repository.SyncProduct(domain.ProductRecord{
		Type:        "category",
		MagentoID:   1,
		DashboardID: categoryData.ID,
	})

	return nil
}

func (c *categoryService) UpdateCategory(consume *domain.Consume) error {
	return nil
}

func (c *categoryService) DeleteCategory(consume *domain.Consume) error {
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
