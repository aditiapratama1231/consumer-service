package controller

import (
	"log"
	"magento-consume-service/database/models"
	"magento-consume-service/domain"

	"github.com/jinzhu/gorm"
)

// Controller structs /
type Controller struct {
	DB      *gorm.DB
	Consume *domain.Consume
}

// MainController function, as main controller, handle incoming data record
func (c *Controller) MainController() {
	var sequenceNumber models.KinesisSequenceNumber
	/*
		check if sequence number already exist in database.
		If true (no error), it was indicate thats data stream already consume.
		So will continue to check next data stream.
	*/
	if err := c.DB.Where("sequence_number = ?", *c.Consume.SequenceNumber).First(&sequenceNumber).Error; err != nil {
		switch service := c.Consume.Data.Head.Service; service {
		case "product":
			c.productManagement()
		case "order":
			c.orderManagement()
		default:
			log.Println("wrong service input")
		}
	}
}
