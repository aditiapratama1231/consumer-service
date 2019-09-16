package domain

type ProductRecord struct {
	Type        string
	DashboardID int
	MagentoID   int
	SKU         string
}

type Product struct {
	ProductSupplierID int
	Name              string
	SKU               string
	Status            int
	Visibility        int
	TypeID            string
	CustomAttributes  []customAttributes
}

type customAttributes struct {
	AttributeCode string
	Value         string
}
