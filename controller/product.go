package controller

import (
	"log"
	"magento-consumer-service/service"
)

func (c *Controller) productManagement() {
	srv := &service.Service{DB: c.DB, Consume: c.Consume}
	switch domain := c.Consume.Data.Head.Domain; domain {
	case "category":
		switch action := c.Consume.Data.Head.Action; action {
		case "create":
			err := srv.CreateCategory()
			if err != nil {
				log.Println(err)
			}
		case "update":
			err := srv.UpdateCategory()
			if err != nil {
				log.Println(err)
			}
		case "delete":
			err := srv.DeleteCategory()
			if err != nil {
				log.Println(err)
			}
		default:
			log.Println("wrong action input")
		}
	default:
		log.Println("wrong domain input")
	}
}
