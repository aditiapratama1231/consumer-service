package service

import (
	"encoding/json"
	"strconv"

	clog "log"

	log "github.com/sirupsen/logrus"

	"magento-consumer-service/config"
	"magento-consumer-service/domain"

	order "magento-consumer-service/order_management"
)

type orderService struct {
	Repository order.OrderRepository
	Request    config.Request
}

func NewOrderService(repo order.OrderRepository, req config.Request) OrderService {
	return &orderService{
		Repository: repo,
		Request:    req,
	}
}

type MagentoResponse struct {
	EntityID int `json:"entity_id"`
}

func (order *orderService) UpdateStatusOrder(consume *domain.Consume) error {
	var (
		magentoResponse MagentoResponse
	)
	payload := consume.Data.Body.Payload
	reqBody, err := json.Marshal(payload)
	if err != nil {
		log.Error("Error encoding product payload " + err.Error())
		return err
	}

	req, err := order.Request.Post("/orders", reqBody)
	resp := req.Response()

	if err != nil || resp.StatusCode != 200 {
		// if get error, show request details
		clog.Printf("%+v", req) // force to show request details in command line.
		log.Printf("%+v", req)
		return err
	}

	req.ToJSON(&magentoResponse)

	_, err = order.Repository.SaveStream(*consume.SequenceNumber)
	if err != nil {
		return err
	}

	dashboardID, err := strconv.Atoi(consume.Data.Head.DashboardID)
	_, err = order.Repository.SyncOrder(domain.OrderRecord{
		Type:        "order",
		DashboardID: dashboardID,
		MagentoID:   magentoResponse.EntityID,
	})

	config.SetAPILogger(req, resp)
	return nil
}

func convertOrderRecord(data interface{}) domain.OrderRecord {
	m := data.(map[string]interface{})
	order := domain.OrderRecord{}
	if order_id, ok := m["order_id"].(int); ok {
		order.OrderID = order_id
	}

	if status, ok := m["status"].(string); ok {
		order.Status = status
	}

	return order
}
