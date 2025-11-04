package main

import (
	"fmt"
	"pos-app-with-fx/models"
	"pos-app-with-fx/services"

	"go.uber.org/fx"
)

func main() {
	// db := &services.InMemoryDB{}
	// logger := &services.ConsoleLogger{}

	// inventory := services.NewInventoryService(db, logger)
	// sales := services.NewSalesService(inventory, logger)

	// inventory.AddItem("Apple", 1.7, 10)
	// total, err := sales.Checkout([]models.ItemSell{{Name: "Apple", Qty: 5}})
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// 	return
	// }
	// fmt.Printf("Total: $%.2f\n", total)

	app := fx.New(
		fx.Provide(
			services.NewInMemoryDB,
			services.NewConsoleLogger,
			services.NewInventoryService,
			services.NewSalesService,
		),
		fx.Invoke(func(inventory *services.InventoryService, sales *services.SalesService) {
			inventory.AddItem("Apple", 2.75, 10)
			total, err := sales.Checkout([]models.ItemSell{{Name: "Apple", Qty: 2}})
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			fmt.Printf("Total $%.2f\n", total)
		}),
	)
	app.Run()
}
