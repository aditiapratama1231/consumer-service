package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"magento-consumer-service/config"
	"magento-consumer-service/consumer"
	_controller "magento-consumer-service/controller"

	_brandService "magento-consumer-service/product_management/brand/service"
	_categoryService "magento-consumer-service/product_management/category/service"
	_productService "magento-consumer-service/product_management/product/service"

	_productRepository "magento-consumer-service/product_management/repository"
)

func main() {
	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	magentoBaseURL := os.Getenv("MAGENTO_BASE_URL")

	// initiate request configuration
	request := config.NewRequest(magentoBaseURL)
	err := request.SetToken()
	if err != nil {
		log.Println("Error getting token : " + err.Error())
	}

	db := config.DBInit()

	// initiate instance
	productRepository := _productRepository.NewProductRepository(db)

	productService := _productService.NewProductService(db, productRepository, request)
	categoryService := _categoryService.NewCategoryService(db, productRepository, request)
	brandService := _brandService.NewBrandService(db, productRepository, request)

	// initiate controller
	controller := _controller.NewController(db,
		categoryService,
		brandService,
		productService)

	errChan := make(chan error)

	// Consumer service
	go func() {
		log.Println("ready to consume data")
		for {
			// initial consumer configuration
			consumer := consumer.NewConsumer(db, controller)
			err := consumer.MainConsumer()
			if err != nil {
				log.Println("Error Get Consumer data", err)
			}
			time.Sleep(30000 * time.Millisecond)
		}
	}()

	log.Fatalln(<-errChan)
}
