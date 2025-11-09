package orders

import (
	"fmt"
	"log"
	"pos/internal/modules/inventory"
)

type orderService struct {
	orderRepo   OrderRepository
	paymentGate PaymentGateway
	invSvc      InventoryServicePort
}

func NewOrderService(or OrderRepository, pg PaymentGateway, is InventoryServicePort) *orderService {
	return &orderService{
		orderRepo:   or,
		paymentGate: pg,
		invSvc:      is,
	}
}

func (s *orderService) CreateOrder(customerID string, itemsToOrder []CreateOrderItem, paymentToken string) (*Order, error) {
	hydratedItems := make([]CartItem, 0)
	stockUpdateItems := make([]inventory.OrderItem, 0)

	log.Println("ORDER Module: Hydrating items and checking stock...")

	for _, item := range itemsToOrder {
		product, err := s.invSvc.GetProduct(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %s not found: %w", item.ProductID, err)
		}

		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("not enough stock for %s. Have %d, need %d", product.Name, product.Stock, item.Quantity)

		}

		hydratedItems = append(hydratedItems, CartItem{
			ID:       product.ID,
			Name:     product.Name,
			Price:    product.Price,
			Quantity: item.Quantity,
		})
		stockUpdateItems = append(stockUpdateItems, inventory.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	order, err := NewOrder(customerID, hydratedItems)
	if err != nil {
		return nil, err
	}
	log.Printf("ORDER Module: Created new order %s with total %.2f", order.ID, order.Total)
	txID, err := s.paymentGate.ProcessPayment(order.Total, paymentToken)
	if err != nil {
		order.Status = StatusFailed
		_ = s.orderRepo.SaveOrder(order)
		return nil, fmt.Errorf("payment failed: %w", err)
	}
	order.Status = StatusPaid
	order.TransactionID = txID
	log.Printf("ORDER Module: Payment successful for order %s (Tx: %s)", order.ID, txID)

	if err := s.invSvc.UpdateStock(stockUpdateItems); err != nil {
		log.Printf("FATAL ERROR: Order %s paid but stock update failed: %v. MANUAL REFUND REQUIRED.", order.ID, err)
		order.Status = StatusFailed
		_ = s.orderRepo.SaveOrder(order)
		return nil, fmt.Errorf("payment OK, but stock update failed: %w. REFUND PENDING.", err)
	}
	log.Printf("ORDER Module: Stock update successful for order %s", order.ID)
	if err := s.orderRepo.SaveOrder(order); err != nil {
		log.Printf("FATAL ERROR: Order %s paid, stock updated, but DB save failed: %v. MANUAL FIX REQUIRED.", order.ID, err)
		return nil, fmt.Errorf("payment and stock OK, but final save failed: %w", err)
	}
	log.Printf("ORDER Module: Successfully processed and saved order %s", order.ID)
	return order, nil
}

func (s *orderService) GetOrderByID(id string) (*Order, error) {
	log.Printf("ORDER Module: Getting order %s", id)
	return s.orderRepo.GetOrderByID(id)
}
