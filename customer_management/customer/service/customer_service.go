package service

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"

	"magento-consumer-service/config"
	customer "magento-consumer-service/customer_management"
	"magento-consumer-service/domain"
)

type customerService struct {
	DB         *gorm.DB
	Request    config.Request
	Repository customer.CustomerRepository
}

func NewCustomerService(db *gorm.DB, repository customer.CustomerRepository, request config.Request) CustomerService {
	return &customerService{
		DB:         db,
		Request:    request,
		Repository: repository,
	}
}

type MagentoResponse struct {
	ID int `json:id`
}

func (c *customerService) CreateCustomer(consume *domain.Consume) error {
	var (
		magentoResponse MagentoResponse
	)

	// convert category payload
	payload := consume.Data.Body.Payload
	fmt.Println("req body", payload)

	reqBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error Encoding customer payload : " + err.Error())
	}

	// POST Data
	req, err := c.Request.Post("/customers", reqBody)
	resp := req.Response()

	if err != nil || resp.StatusCode != 200 {
		// if get error, show request details
		log.Printf("%+v", req)
		log.Println("Error SetUp API call : ", err)
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
	_, err = c.Repository.SyncCustomer(domain.CustomerRecord{
		Type:        "customer",
		MagentoID:   magentoResponse.ID,
		DashboardID: dashboardID,
	})

	return nil

}

//UpdateCustomer /
func (c *customerService) UpdateCustomer(consume *domain.Consume) error {
	var magentoResponse MagentoResponse

	payload := consume.Data.Body.Payload
	reqBody, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error encoding customer payload " + err.Error())
		return err
	}

	dashboardID, err := strconv.Atoi(consume.Data.Head.Dashboard)
	ctgry, err := c.Repository.GetMagentoID("customer", dashboardID)
	if err != nil {
		log.Println("Error get magento id from database : " + err.Error())
		return err
	}

	endpoint := "/customers/" + strconv.Itoa(ctgry.MagentoID)
	req, err := c.Request.Put(endpoint, reqBody)
	resp := req.Response()
	if err != nil || resp.StatusCode != 200 {
		// if get error, show request details
		log.Printf("%+v", req)
		return err
	}

	req.ToJSON(&magentoResponse)

	_, err = c.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		log.Println("Error save stream to database : " + err.Error())
		return err
	}

	_, err = c.Repository.SyncCustomer(domain.CustomerRecord{
		Type:        "customer",
		MagentoID:   magentoResponse.ID,
		DashboardID: dashboardID,
	})

	if err != nil {
		log.Println("Error sync customer to database : " + err.Error())
		return err
	}

	return nil

}

func (customer *customerService) DeleteCustomer(consume *domain.Consume) error {

	return nil
}
