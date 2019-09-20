package models

// ProductRecord struct models, collection of record product management
type OrderRecord struct {
	Model
	Type        string `gorm:"column:type"`
	DashboardID int    `gorm:"column:dashboard_id"`
	MagentoID   int    `gorm:"column:magento_id"`
}
