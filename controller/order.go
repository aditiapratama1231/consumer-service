package controller

import (
	"magento-consumer-service/domain"

	log "github.com/sirupsen/logrus"
)

func (c *controller) OrderManagement(consume *domain.Consume) {
	switch domain := consume.Data.Head.Domain; domain {
	case "order":
		switch action := consume.Data.Head.Action; action {
		case "update":
			err := c.OrderService.UpdateStatusOrder(consume)
			if err != nil {
				log.Error(err)
			}
		default:
			log.Fatal("wrong action input")
		}
	}
	log.Println("order")
}
