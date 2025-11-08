package db

import (
	"fmt"
	"log"
	"pos/internal/modules/inventory"
	"pos/internal/modules/orders"
	"sync"
)

type InMemoryAdapter struct {
	products map[string]*inventory.Product
	orders   map[string]*orders.Order
	mu       sync.RWMutex
}

func NewInMemoryAdapter() *InMemoryAdapter {
	return &InMemoryAdapter{
		products: make(map[string]*inventory.Product),
		orders:   make(map[string]*orders.Order),
	}
}

func (a *InMemoryAdapter) SaveProduct(product *inventory.Product) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.products[product.ID] = product
	log.Printf("DB Adapter: Saved product %s", product.Name)
	return nil
}

func (a *InMemoryAdapter) GetProduct(id string) (*inventory.Product, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	product, ok := a.products[id]
	if !ok {
		return nil, fmt.Errorf("product %s not found in DB", id)
	}
	return product, nil
}

func (a *InMemoryAdapter) UpdateProductStock(id string, newStock int) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	product, ok := a.products[id]
	if !ok {
		return fmt.Errorf("product %s not found in DB", id)
	}
	product.Stock = newStock
	log.Printf("DB Adapter: Updated stock for %s to %d", product.Name, newStock)
	return nil
}

func (a *InMemoryAdapter) SaveOrder(order *orders.Order) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.orders[order.ID] = order
	log.Printf("DB Adapter: Saved order %s (Status: %s)", order.ID, order.Status)
	return nil
}

func (a *InMemoryAdapter) GetOrderByID(id string) (*orders.Order, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	order, ok := a.orders[id]
	if !ok {
		return nil, fmt.Errorf("order %s not found in DB", id)
	}
	return order, nil
}
