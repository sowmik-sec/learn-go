package main

import (
	"log"
	"pos/internal/adapters/driven/db"
	"pos/internal/adapters/driven/payment"
	myHttp "pos/internal/adapters/driving/http"
	"pos/internal/modules/inventory"
	"pos/internal/modules/orders"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func main() {
	log.Println("--- Starting Hexagonal Modular Monolith POS ---")
	app := fx.New(

		fx.Provide(
			fx.Annotate(db.NewInMemoryAdapter, fx.As(new(inventory.ProductRepository)), fx.As(new(orders.OrderRepository))),
		),
		fx.Provide(
			fx.Annotate(payment.NewMockPaymentGateway, fx.As(new(orders.PaymentGateway))),
		),

		fx.Provide(fx.Annotate(inventory.NewInventoryService, fx.As(new(inventory.InventoryService)), fx.As(new(orders.InventoryServicePort)))),
		fx.Provide(fx.Annotate(orders.NewOrderService, fx.As(new(orders.OrderService)))),

		fx.Provide(myHttp.NewHTTPHandler),
		fx.Provide(func() *gin.Engine { return gin.Default() }),

		fx.Invoke(startServer),
	)

	app.Run()
}

func startServer(handler *myHttp.HTTPHandler, router *gin.Engine) {
	handler.RegisterRoutes(router)

	log.Println("Server starting on :8080...")
	log.Println("Endpoints:")
	log.Println("  POST /product/add    (See AddProductRequest)")
	log.Println("  GET  /product/get?id= (Returns Product)")
	log.Println("  GET  /product/all    (Returns []Product)")
	log.Println("  POST /order/create   (See CreateOrderRequest)")
	log.Println("  GET  /order/get?id=   (Returns Order)")

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
