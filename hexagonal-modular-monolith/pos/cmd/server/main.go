package main

import (
	"log"
	http "net/http"
	"pos/internal/adapters/driven/db"
	"pos/internal/adapters/driven/payment"
	myHttp "pos/internal/adapters/driving/http"
	"pos/internal/modules/inventory"
	"pos/internal/modules/orders"
)

func main() {
	log.Println("--- Starting Hexagonal Modular Monolith POS ---")
	dbAdapter := db.NewInMemoryAdapter()
	paymentGate := payment.NewMockPaymentGateway()

	inventoryService := inventory.NewInventoryService(dbAdapter)

	orderService := orders.NewOrderService(dbAdapter, paymentGate, inventoryService)

	httpHandler := myHttp.NewHTTPHandler(orderService, inventoryService)

	mux := http.NewServeMux()
	httpHandler.RegisterRoutes(mux)

	log.Println("Server starting on :8080...")
	log.Println("Endpoints:")
	log.Println("  POST /product/add    (See AddProductRequest)")
	log.Println("  GET  /product/get?id= (Returns Product)")
	log.Println("  POST /order/create   (See CreateOrderRequest)")
	log.Println("  GET  /order/get?id=   (Returns Order)")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
