package inventory

import (
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type inventoryService struct {
	repo ProductRepository
}

func NewInventoryService(repo ProductRepository) *inventoryService {
	return &inventoryService{repo: repo}
}

func (s *inventoryService) AddProduct(name string, price float64, stock int) (*Product, error) {
	if name == "" || price <= 0 || stock < 0 {
		return nil, errors.New("invalid product data")
	}
	// Check for duplicate product name
	_, err := s.repo.GetProductByName(name)
	if err == nil {
		return nil, fmt.Errorf("product with name '%s' already exists", name)
	}
	product := &Product{
		ID:    uuid.NewString(),
		Name:  name,
		Price: price,
		Stock: stock,
	}
	if err := s.repo.SaveProduct(product); err != nil {
		return nil, err
	}
	log.Printf("INVENTORY Module: Added product %s (ID: %s)", name, product.ID)
	return product, nil
}

func (s *inventoryService) GetProduct(id string) (*Product, error) {
	log.Printf("INVENTORY Module: Getting Product %s", id)
	return s.repo.GetProduct(id)
}

func (s *inventoryService) GetAllProducts() ([]*Product, error) {
	log.Printf("INVENTORY Module: Getting all products")
	return s.repo.GetAllProducts()
}

func (s *inventoryService) UpdateStock(items []OrderItem) error {
	log.Printf("INVENTORY Module: Updating stock for %d items", len(items))
	for _, item := range items {
		product, err := s.repo.GetProduct(item.ProductID)
		if err != nil {
			return fmt.Errorf("could not get product %s: %w", item.ProductID, err)
		}
		if product.Stock < item.Quantity {
			return fmt.Errorf("not enough stock for %s. have %d, need %d", product.Name, product.Stock, item.Quantity)
		}

		newStock := product.Stock - item.Quantity
		if err := s.repo.UpdateProductStock(product.ID, newStock); err != nil {
			return fmt.Errorf("failed to update stock for %s: %w", product.ID, err)
		}
		log.Printf("INVENTORY Module: Updated stock for %s. Old: %d, New %d", product.Name, product.Stock, newStock)

	}
	return nil
}
