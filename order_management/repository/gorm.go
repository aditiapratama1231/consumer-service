package repository

import (
	"magento-consumer-service/database/models"
	"magento-consumer-service/domain"
	order "magento-consumer-service/order_management"

	"github.com/jinzhu/gorm"
)

type orderRepository struct {
	DB *gorm.DB
}

func NewOrderRepository(db *gorm.DB) order.OrderRepository {
	return &orderRepository{
		DB: db,
	}
}

func (o *orderRepository) SyncOrder(data domain.OrderRecord) (interface{}, error) {
	var order models.OrderRecord
	tx := o.DB.Begin()

	order.Type = data.Type
	order.DashboardID = data.DashboardID
	order.MagentoID = data.MagentoID

	if tx.Where("type = ?", data.Type).
		Where("magento_id = ?", data.MagentoID).
		Where("dashboard_id = ?", data.DashboardID).
		Find(&order).RecordNotFound() {

		err := tx.Create(&order).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		tx.Commit()
		return order, nil
	}

	err := tx.Save(&order).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return order, nil
}

func (o *orderRepository) GetMagentoID(tp string, dashboardId int) (models.OrderRecord, error) {
	var order models.OrderRecord

	o.DB.Where("dashboard_id = ?", dashboardId).
		Where("type = ?", tp).
		First(&order)

	return order, nil
}

func (o *orderRepository) SaveStream(sequenceNumber string) (interface{}, error) {
	var kinesis models.KinesisSequenceNumber

	kinesis.SequenceNumber = sequenceNumber

	tx := o.DB.Begin()

	err := tx.Create(&kinesis).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return kinesis, nil
}

func (o *orderRepository) DeleteRecord(tp string, dashboardID int) error {
	var order models.OrderRecord
	tx := o.DB.Begin()

	err := tx.Where("dashboard_id=?", dashboardID).
		Where("type = ?", tp).
		First(&order).Delete(&order).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil

}
