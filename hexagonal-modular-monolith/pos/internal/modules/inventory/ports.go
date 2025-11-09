package inventory

type InventoryService interface {
	AddProduct(name string, price float64, stock int) (*Product, error)
	GetProduct(id string) (*Product, error)
	GetAllProducts() ([]*Product, error)
	UpdateStock(items []OrderItem) error
}

type ProductRepository interface {
	SaveProduct(product *Product) error
	GetProduct(id string) (*Product, error)
	GetAllProducts() ([]*Product, error)
	GetProductByName(name string) (*Product, error)
	UpdateProductStock(id string, newStock int) error
}

type OrderItem struct {
	ProductID string
	Quantity  int
}
