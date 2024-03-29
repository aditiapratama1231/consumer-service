package service

import (
	"encoding/json"
	clog "log"
	"strconv"

	log "github.com/sirupsen/logrus"

	"magento-consumer-service/config"
	"magento-consumer-service/domain"
	product "magento-consumer-service/product_management"
)

type productService struct {
	Repository product.ProductRepository
	Request    config.Request
}

func NewProductService(repository product.ProductRepository, request config.Request) ProductService {
	return &productService{
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
		log.Error("Error encoding product payload " + err.Error())
	}

	req, err := product.Request.Post("/products", reqBody)
	resp := req.Response()

	if err != nil || resp.StatusCode != 200 {
		// if get error, show request details
		clog.Printf("%+v", req) // force to show request details in command line.
		log.Errorf("%+v", req)
		log.Error("Error SetUp API call : ", err)
		return err
	}

	req.ToJSON(&magentoResponse)

	_, err = product.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		return err
	}

	productPayload := convertProductRecord(payload["product"])
	dashboardID, err := strconv.Atoi(consume.Data.Head.DashboardID)
	_, err = product.Repository.SyncProduct(domain.ProductRecord{
		Type:        "product",
		MagentoID:   magentoResponse.ID,
		DashboardID: dashboardID,
		SKU:         productPayload.SKU,
	})

	if err != nil {
		log.Error("Error sync product to database : " + err.Error())
		return err
	}

	config.SetAPILogger(req, resp)
	return nil
}

func (product *productService) UpdateProduct(consume *domain.Consume) error {
	var (
		magentoResponse MagentoResponse
	)

	payload := consume.Data.Body.Payload
	reqBody, err := json.Marshal(payload)
	if err != nil {
		log.Error("Error encoding product payload " + err.Error())
		return err
	}

	dashboardID, err := strconv.Atoi(consume.Data.Head.DashboardID)
	prd, err := product.Repository.GetMagentoID("product", dashboardID)
	if err != nil {
		log.Error("Error get magento id from database : " + err.Error())
		return err
	}

	endpoint := "/products/" + prd.SKU
	req, err := product.Request.Put(endpoint, reqBody)
	resp := req.Response()

	if err != nil || resp.StatusCode != 200 {
		// if get error, show request details
		clog.Printf("%+v", req) // force to show request details in command line.
		log.Printf("%+v", req)
		return err
	}

	req.ToJSON(&magentoResponse)
	_, err = product.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		log.Error("Error save stream to database : " + err.Error())
		return err
	}

	_, err = product.Repository.SyncProduct(domain.ProductRecord{
		Type:        "product",
		MagentoID:   magentoResponse.ID,
		DashboardID: dashboardID,
		SKU:         prd.SKU,
	})

	if err != nil {
		log.Error("Error sync product to database : " + err.Error())
		return err
	}

	config.SetAPILogger(req, resp)
	return nil
}

func (product *productService) DeleteProduct(consume *domain.Consume) error {
	var (
		magentoResponse MagentoResponse
	)

	dashboardID, err := strconv.Atoi(consume.Data.Head.DashboardID)
	prd, err := product.Repository.GetMagentoID("product", dashboardID)
	if err != nil {
		log.Error("Error get magento id from database : " + err.Error())
		return err
	}

	endpoint := "/products/" + prd.SKU
	req, err := product.Request.Delete(endpoint)
	resp := req.Response()

	if err != nil || resp.StatusCode != 200 {
		// if get error, show request details
		clog.Printf("%+v", req) // force to show request details in command line.
		log.Printf("%+v", req)
		log.Error("Error SetUp API call : ", err)
		return err
	}

	req.ToJSON(&magentoResponse)
	_, err = product.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		log.Error("Error save stream to database : " + err.Error())
		return err
	}

	err = product.Repository.DeleteRecord(consume.Data.Head.Domain, dashboardID)
	if err != nil {
		log.Println("Error delete record in database: " + err.Error())
		return err
	}

	config.SetAPILogger(req, resp)
	return nil
}

func convertProductRecord(data interface{}) domain.ProductRecord {
	m := data.(map[string]interface{})
	product := domain.ProductRecord{}
	if sku, ok := m["sku"].(string); ok {
		product.SKU = sku
	}

	return product
}
