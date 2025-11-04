package models

type Item struct {
	Name  string
	Price float64
	Stock int
}

type ItemSell struct {
	Name string
	Qty  int
}
