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
	product "magento-consumer-service/product_management"
	_categoryService "magento-consumer-service/product_management/category/service"
)

func main() {
	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	magentoBaseURL := os.Getenv("MAGENTO_BASE_URL")

	db := config.DBInit()
	request := config.NewRequest(magentoBaseURL)

	// initial
	productRepository := product.NewProductRepository(db)

	categoryService := _categoryService.NewCategoryService(db, productRepository, request)
	controller := _controller.NewController(db, categoryService)

	errChan := make(chan error)

	// Consumer service
	go func() {
		log.Println("ready to consume data")
		for {
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
