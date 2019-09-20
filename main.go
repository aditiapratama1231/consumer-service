package main

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"magento-consumer-service/config"
	"magento-consumer-service/consumer"
	_controller "magento-consumer-service/controller"

	_customerService "magento-consumer-service/customer_management/customer/service"
	_orderService "magento-consumer-service/order_management/order/service"
	_brandService "magento-consumer-service/product_management/brand/service"
	_categoryService "magento-consumer-service/product_management/category/service"
	_productService "magento-consumer-service/product_management/product/service"

	_customerRepository "magento-consumer-service/customer_management/repository"
	_orderRepository "magento-consumer-service/order_management/repository"
	_productRepository "magento-consumer-service/product_management/repository"
)

func main() {
	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	// Setup log
	config.LogFormatter()

	magentoBaseURL := os.Getenv("MAGENTO_BASE_URL")

	// initiate request configuration
	request := config.NewRequest(magentoBaseURL)
	err := request.SetToken()
	if err != nil {
		log.Info("Error getting token : " + err.Error())
	}

	db := config.DBInit()

	// initiate instance
	// initiate repository
	productRepository := _productRepository.NewProductRepository(db)
	customerRepository := _customerRepository.NewCustomerRepository(db)
	orderRepository := _orderRepository.NewOrderRepository(db)

	// product management entity
	productService := _productService.NewProductService(db, productRepository, request)
	categoryService := _categoryService.NewCategoryService(db, productRepository, request)
	brandService := _brandService.NewBrandService(db, productRepository, request)
	orderService := _orderService.NewOrderService(orderRepository, request)

	// customer entity
	customerService := _customerService.NewCustomerService(db, customerRepository, request)

	// initiate controller
	controller := _controller.NewController(db,
		categoryService,
		brandService,
		productService,
		customerService,
		orderService)

	errChan := make(chan error)

	// Consumer service
	go func() {
		log.Info("ready to consume data")
		for {
			// initial consumer configuration
			consumer := consumer.NewConsumer(db, controller)
			err := consumer.MainConsumer()
			if err != nil {
				log.Error("Error Get Consumer data", err)
			}
			time.Sleep(30000 * time.Millisecond)
		}
	}()

	log.Fatal(<-errChan)
}
