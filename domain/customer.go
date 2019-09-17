package domain

// Customer struct /
type Customer struct {
	CustomerAttr
	Password string
}

//CustomerAttr /
type CustomerAttr struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	Addresses []*Address
}

//Address /
type Address struct {
	DefaultShipping bool
	DefaultBilling  bool
	FirstName       string
	LastName        string
	Region          *Region
	PostCode        string
	Street          []string
	Telephone       string
	CountryID       string
}

//Region /
type Region struct {
	RegionCode string
	Region     string
	RegionID   int
}

type CustomerRecord struct {
	Type        string
	DashboardID int
	MagentoID   int
}
