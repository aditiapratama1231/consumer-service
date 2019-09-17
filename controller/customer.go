package controller

import (
	"log"
	"magento-consumer-service/domain"
)

func (c *controller) CustomerManagement(consume *domain.Consume) {
	switch domain := consume.Data.Head.Domain; domain {
	case "customer":
		switch action := consume.Data.Head.Action; action {
		case "create":
			err := c.CustomerService.CreateCustomer(consume)
			if err != nil {
				log.Println("error create", err)
			}
		case "update":
			err := c.CustomerService.UpdateCustomer(consume)
			if err != nil {
				log.Println(err)
			}
		case "delete":
			err := c.CustomerService.DeleteCustomer(consume)
			if err != nil {
				log.Println(err)
			}
		default:
			log.Println("wrong action input")
		}
	}
}
