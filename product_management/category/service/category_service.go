package service

import (
	"encoding/json"
	"errors"
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
	if consume.Data.Body.Payload["category"] != nil {
		categoryData = convertCategory(consume.Data.Body.Payload["category"])
	} else {
		return errors.New("wrong payload")
	}

	reqBody, err := json.Marshal(consume.Data)
	if err != nil {
		return err
	}

	// POST Data
	req, err := c.Request.Post("/products", reqBody)
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
	_, err = c.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		return err
	}

	// save to record
	productRecord := domain.ProductRecord{
		Type:        "category",
		MagentoID:   1,
		DashboardID: categoryData.ID,
	}
	_, err = c.Repository.SyncProduct(productRecord)

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
