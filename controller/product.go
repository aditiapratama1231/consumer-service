package controller

import (
	"log"
	"magento-consumer-service/domain"
)

func (c *controller) ProductManagement(consume *domain.Consume) {
	switch domain := consume.Data.Head.Domain; domain {

	case "category":
		switch action := consume.Data.Head.Action; action {
		case "create":
			err := c.CategoryService.CreateCategory(consume)
			if err != nil {
				log.Println(err)
			}
		case "update":
			err := c.CategoryService.UpdateCategory(consume)
			if err != nil {
				log.Println(err)
			}
		case "delete":
			err := c.CategoryService.DeleteCategory(consume)
			if err != nil {
				log.Println(err)
			}
		default:
			log.Println("wrong action input")
		}

	case "brand":
		switch action := consume.Data.Head.Action; action {
		case "create":
			err := c.BrandService.CreateBrand(consume)
			if err != nil {
				log.Println(err)
			}
		case "update":
			err := c.BrandService.UpdateBrand(consume)
			if err != nil {
				log.Println(err)
			}
		case "delete":
			err := c.BrandService.DeleteBrand(consume)
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
