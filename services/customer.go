package services

type Customer struct {
	ID   int
	Name string
}

func GetCustomers() []Customer {
	return []Customer{
		{ID: 1, Name: "ABC Traders"},
		{ID: 2, Name: "XYZ Enterprise"},
		{ID: 3, Name: "Demo Company"},
	}
}
