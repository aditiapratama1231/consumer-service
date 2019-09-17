package models

// CustomerRecord struct models, collection of record product management
type CustomerRecord struct {
	Model
	Type        string `gorm:"column:type"`
	DashboardID int    `gorm:"column:dashboard_id"`
	MagentoID   int    `gorm:"column:magento_id"`
}
