package repository

import (
	customer "magento-consumer-service/customer_management"
	"magento-consumer-service/database/models"
	"magento-consumer-service/domain"

	"github.com/jinzhu/gorm"
)

type customerRepository struct {
	DB *gorm.DB
}

//NewCustomerRepository /
func NewCustomerRepository(DB *gorm.DB) customer.CustomerRepository {
	return &customerRepository{
		DB: DB,
	}
}

func (c *customerRepository) SyncCustomer(data domain.CustomerRecord) (interface{}, error) {
	var customer models.CustomerRecord
	tx := c.DB.Begin()

	customer.Type = data.Type
	customer.DashboardID = data.DashboardID
	customer.MagentoID = data.MagentoID

	if tx.Where("type =?", data.Type).
		Where("magento_id=?", data.MagentoID).
		Where("dashboard_id=?", data.DashboardID).
		Find(&customer).RecordNotFound() {
		err := tx.Create(&customer).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		tx.Commit()
		return customer, nil
	}
	err := tx.Save(&customer).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return customer, nil

}

func (c *customerRepository) GetMagentoID(tp string, dashboardID int) (models.CustomerRecord, error) {
	var customer models.CustomerRecord
	c.DB.Where("dashboard_id = ?", dashboardID).
		Where("type = ?", tp).
		First(&customer)

	return customer, nil
}

func (c *customerRepository) ShowCustomerID(data domain.CustomerRecord) (interface{}, error) {
	var customer models.CustomerRecord

	c.DB.Where("type = ?", data.Type).Where("magento_id = ?", data.MagentoID).Where("dashboard_id = ?", data.DashboardID).Find(&customer)

	return customer, nil
}

func (c *customerRepository) SaveStream(sequenceNumber string) (interface{}, error) {
	var kinesis models.KinesisSequenceNumber

	kinesis.SequenceNumber = sequenceNumber

	tx := c.DB.Begin()

	err := tx.Create(&kinesis).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	return kinesis, nil
}
