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

type MagentoResponse struct {
	ID int `json:"id"`
}

func (product *productService) CreateProduct(consume *domain.Consume) error {
	var (
		magentoResponse MagentoResponse
	)

	payload := consume.Data.Body.Payload
	reqBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error encoding product payload " + err.Error())
	}

	req, err := product.Request.Post("/products", reqBody)
	if err != nil {
		// if get error, show request details
		log.Printf("%+v", req)
		log.Println("Error SetUp API call : ", err.Error())
		return err
	}

	req.ToJSON(&magentoResponse)

	_, err = product.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		return err
	}

	productPayload := convertProductRecord(payload["product"])
	dashboardID, err := strconv.Atoi(consume.Data.Head.Dashboard)
	_, err = product.Repository.SyncProduct(domain.ProductRecord{
		Type:        "product",
		MagentoID:   magentoResponse.ID,
		DashboardID: dashboardID,
		SKU:         productPayload.SKU,
	})

	if err != nil {
		log.Println("Error sync product to database : " + err.Error())
		return err
	}

	return nil
}

func (product *productService) UpdateProduct(consume *domain.Consume) error {
	var (
		magentoResponse MagentoResponse
	)

	payload := consume.Data.Body.Payload
	reqBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error encoding product payload " + err.Error())
		return err
	}

	dashboardID, err := strconv.Atoi(consume.Data.Head.Dashboard)
	prd, err := product.Repository.GetMagentoID(dashboardID)
	if err != nil {
		log.Println("Error get magento id from database : " + err.Error())
		return err
	}

	endpoint := "/products/" + prd.SKU
	req, err := product.Request.Put(endpoint, reqBody)
	if err != nil {
		log.Printf("%+v", req)
		return err
	}

	req.ToJSON(&magentoResponse)
	_, err = product.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		log.Println("Error save stream to database : " + err.Error())
		return err
	}

	_, err = product.Repository.SyncProduct(domain.ProductRecord{
		Type:        "product",
		MagentoID:   magentoResponse.ID,
		DashboardID: dashboardID,
		SKU:         prd.SKU,
	})

	if err != nil {
		log.Println("Error sync product to database : " + err.Error())
		return err
	}

	return nil
}

func (product *productService) DeleteProduct(consume *domain.Consume) error {
	var (
		magentoResponse MagentoResponse
	)

	dashboardID, err := strconv.Atoi(consume.Data.Head.Dashboard)
	prd, err := product.Repository.GetMagentoID(dashboardID)
	if err != nil {
		log.Println("Error get magento id from database : " + err.Error())
		return err
	}

	endpoint := "/products/" + prd.SKU
	req, err := product.Request.Delete(endpoint)

	if err != nil {
		log.Printf("%+v", req)
		return err
	}

	req.ToJSON(&magentoResponse)
	_, err = product.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		log.Println("Error save stream to database : " + err.Error())
		return err
	}

	return nil
}

func convertProductRecord(data interface{}) domain.ProductRecord {
	m := data.(map[string]interface{})
	log.Println(m["sku"])
	product := domain.ProductRecord{}
	if sku, ok := m["sku"].(string); ok {
		product.SKU = sku
	}

	return product
}
