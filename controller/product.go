package controller

import (
	"magento-consumer-service/domain"

	log "github.com/sirupsen/logrus"
)

func (c *controller) ProductManagement(consume *domain.Consume) {
	switch domain := consume.Data.Head.Domain; domain {
	case "product":
		switch action := consume.Data.Head.Action; action {
		case "create":
			err := c.ProductService.CreateProduct(consume)
			if err != nil {
				log.Error(err)
			}
		case "update":
			err := c.ProductService.UpdateProduct(consume)
			if err != nil {
				log.Error(err)
			}
		case "delete":
			err := c.ProductService.DeleteProduct(consume)
			if err != nil {
				log.Error(err)
			}
		default:
			log.Fatal("wrong action input")
		}

	case "category":
		switch action := consume.Data.Head.Action; action {
		case "create":
			err := c.CategoryService.CreateCategory(consume)
			if err != nil {
				log.Error(err)
			}
		case "update":
			err := c.CategoryService.UpdateCategory(consume)
			if err != nil {
				log.Error(err)
			}
		case "delete":
			err := c.CategoryService.DeleteCategory(consume)
			if err != nil {
				log.Error(err)
			}
		default:
			log.Println("wrong action input")
		}

	case "brand":
		switch action := consume.Data.Head.Action; action {
		case "create":
			err := c.BrandService.CreateBrand(consume)
			if err != nil {
				log.Error(err)
			}
		case "update":
			err := c.BrandService.UpdateBrand(consume)
			if err != nil {
				log.Error(err)
			}
		case "delete":
			err := c.BrandService.DeleteBrand(consume)
			if err != nil {
				log.Error(err)
			}
		default:
			log.Fatal("wrong action input")
		}
	default:
		log.Fatal("wrong domain input")
	}
}
