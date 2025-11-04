package main

import (
	"fmt"
	"pos-app/models"
	"pos-app/services"
)

// type Logger interface {
// 	Log(msg string)
// }

// type ConsoleLogger struct{}

// func (cl *ConsoleLogger) Log(msg string) {
// 	fmt.Println(msg)
// }

// type MockLogger struct {
// 	logs []string
// }

// func (ml *MockLogger) Log(msg string) {
// 	ml.logs = append(ml.logs, msg)
// }

// type InventoryService struct {
// 	logger Logger
// }

// func NewInventoryService(logger Logger) *InventoryService {
// 	return &InventoryService{
// 		logger: logger,
// 	}
// }

// func (is *InventoryService) AddItem(name string) {
// 	is.logger.Log("Added item: " + name)
// }

// func main() {
// 	realLogger := &ConsoleLogger{}
// 	realService := NewInventoryService(realLogger)
// 	realService.AddItem("Apple")

// 	mockLogger := &MockLogger{}
// 	mockService := NewInventoryService(mockLogger)
// 	mockService.AddItem("Banana")
// 	mockService.AddItem("Orange")

// 	fmt.Println("Mock logs: ", mockLogger.logs)
// }

func main() {
	db := &services.InMemoryDB{}
	logger := &services.ConsoleLogger{}

	inventory := services.NewInventoryService(db, logger)
	sales := services.NewSalesService(inventory, logger)

	inventory.AddItem("Apple", 1.7, 10)
	total, err := sales.Checkout([]models.ItemSell{{Name: "Apple", Qty: 5}})
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf("Total: $%.2f\n", total)
}
