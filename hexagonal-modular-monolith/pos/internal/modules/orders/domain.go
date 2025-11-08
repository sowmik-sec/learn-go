package orders

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusPending OrderStatus = "PENDING"
	StatusPaid    OrderStatus = "PAID"
	StatusFailed  OrderStatus = "FAILED"
)

type CartItem struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

type Order struct {
	ID            string      `json:"id"`
	CustomerID    string      `json:"customer_id"`
	Items         []CartItem  `json:"items"`
	Total         float64     `json:"total"`
	Status        OrderStatus `json:"status"`
	TransactionID string      `json:"transaction_id,omitempty"`
	CreatedAt     time.Time   `json:"created_at"`
}

func NewOrder(customerID string, items []CartItem) (*Order, error) {
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}
	var total float64
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return &Order{
		ID:         uuid.NewString(),
		CustomerID: customerID,
		Items:      items,
		Total:      total,
		Status:     StatusPending,
		CreatedAt:  time.Now().UTC(),
	}, nil
}
