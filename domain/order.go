package domain

type OrderRecord struct {
	Type        string
	DashboardID int
	MagentoID   int
	OrderID     int
	Status      string
}
