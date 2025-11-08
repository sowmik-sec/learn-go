package orders

import (
	"pos/internal/modules/inventory"
)

type OrderService interface {
	CreateOrder(customerID string, itemsToOrder []CreateOrderItem, paymentToken string) (*Order, error)
	GetOrderByID(id string) (*Order, error)
}

type CreateOrderItem struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type OrderRepository interface {
	SaveOrder(order *Order) error
	GetOrderByID(id string) (*Order, error)
}

type PaymentGateway interface {
	ProcessPayment(amount float64, token string) (transactionID string, err error)
}

type InventoryServicePort interface {
	GetProduct(id string) (*inventory.Product, error)
	UpdateStock(items []inventory.OrderItem) error
}
