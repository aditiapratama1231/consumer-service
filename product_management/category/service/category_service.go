package service

import (
	"encoding/json"
	clog "log"
	"strconv"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

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

type MagentoResponse struct {
	ID int `json:"id"`
}

func (c *categoryService) CreateCategory(consume *domain.Consume) error {
	var (
		magentoResponse MagentoResponse
	)

	// convert category payload
	payload := consume.Data.Body.Payload
	reqBody, err := json.Marshal(payload)
	if err != nil {
		log.Error("Error Encoding category payload : " + err.Error())
	}

	// POST Data
	req, err := c.Request.Post("/categories", reqBody)
	resp := req.Response()

	if err != nil || resp.StatusCode != 200 {
		// if get error, show request details
		clog.Printf("%+v", req) // force to show request details in command line.
		log.Printf("%+v", req)
		log.Error("Error SetUp API call : ", err)
		return err
	}

	req.ToJSON(&magentoResponse)

	// if POST success, safe data to db
	_, err = c.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		return err
	}

	// save to record
	dashboardID, err := strconv.Atoi(consume.Data.Head.Dashboard)
	_, err = c.Repository.SyncProduct(domain.ProductRecord{
		Type:        "category",
		MagentoID:   magentoResponse.ID,
		DashboardID: dashboardID,
	})

	config.SetAPILogger(req, resp)
	return nil
}

func (c *categoryService) UpdateCategory(consume *domain.Consume) error {
	var (
		magentoResponse MagentoResponse
	)

	payload := consume.Data.Body.Payload
	reqBody, err := json.Marshal(payload)
	if err != nil {
		log.Error("Error encoding product payload " + err.Error())
		return err
	}

	dashboardID, err := strconv.Atoi(consume.Data.Head.Dashboard)
	ctgry, err := c.Repository.GetMagentoID("category", dashboardID)
	if err != nil {
		log.Error("Error get magento id from database : " + err.Error())
		return err
	}

	endpoint := "/categories/" + strconv.Itoa(ctgry.MagentoID)
	req, err := c.Request.Put(endpoint, reqBody)
	resp := req.Response()
	if err != nil || resp.StatusCode != 200 {
		// if get error, show request details
		clog.Printf("%+v", req) // force to show request details in command line.
		log.Printf("%+v", req)
		return err
	}

	req.ToJSON(&magentoResponse)

	_, err = c.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		log.Error("Error save stream to database : " + err.Error())
		return err
	}

	_, err = c.Repository.SyncProduct(domain.ProductRecord{
		Type:        "category",
		MagentoID:   magentoResponse.ID,
		DashboardID: dashboardID,
	})

	if err != nil {
		log.Error("Error sync product to database : " + err.Error())
		return err
	}

	config.SetAPILogger(req, resp)
	return nil
}

func (c *categoryService) DeleteCategory(consume *domain.Consume) error {
	var (
		magentoResponse MagentoResponse
	)

	dashboardID, err := strconv.Atoi(consume.Data.Head.Dashboard)
	ctgry, err := c.Repository.GetMagentoID("category", dashboardID)
	if err != nil {
		log.Error("Error get magento id from database : " + err.Error())
		return err
	}

	endpoint := "/categories/" + strconv.Itoa(ctgry.MagentoID)
	req, err := c.Request.Delete(endpoint)
	resp := req.Response()
	if err != nil || resp.StatusCode != 200 {
		// if get error, show request details
		clog.Printf("%+v", req) // force to show request details in command line.
		log.Printf("%+v", req)
		log.Error("Error SetUp API call : ", err)
		return err
	}

	req.ToJSON(&magentoResponse)
	_, err = c.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		log.Error("Error save stream to database : " + err.Error())
		return err
	}
	err = c.Repository.DeleteRecord(consume.Data.Head.Domain, dashboardID)
	if err != nil {
		log.Println("Error delete record in database: " + err.Error())
		return err
	}
	config.SetAPILogger(req, resp)
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
