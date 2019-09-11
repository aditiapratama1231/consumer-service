package main

import (
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"

	"magento-consume-service/config"
	"magento-consume-service/consumer"
)

func main() {
	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}

	db := config.DBInit()

	errChan := make(chan error)

	// Consumer service
	go func() {
		log.Println("ready to consume data")
		for {
			err := consumer.Consumer(db)
			if err != nil {
				log.Println(err)
			}
			time.Sleep(30000 * time.Millisecond)
		}
	}()

	log.Fatalln(<-errChan)
}
