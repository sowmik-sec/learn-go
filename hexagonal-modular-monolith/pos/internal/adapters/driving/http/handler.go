package http

import (
	"log"
	"net/http"
	"pos/internal/modules/inventory"
	"pos/internal/modules/orders"
	"strings"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	orderSvc     orders.OrderService
	inventorySvc inventory.InventoryService
}

func NewHTTPHandler(os orders.OrderService, is inventory.InventoryService) *HTTPHandler {
	return &HTTPHandler{
		orderSvc:     os,
		inventorySvc: is,
	}
}

func (h *HTTPHandler) RegisterRoutes(gin *gin.Engine) {
	gin.POST("/order/create", h.handleCreateOrder)
	gin.GET("/order/get", h.handleGetOrder)
	gin.POST("/product/add", h.handleAddProduct)
	gin.GET("/product/get", h.handleGetProduct)
	gin.GET("/product/all", h.handleGetAllProducts)
}

type CreateOrderRequest struct {
	CustomerID   string                   `json:"customer_id"`
	Items        []orders.CreateOrderItem `json:"items"`
	PaymentToken string                   `json:"payment_token"`
}

type AddProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

func (h *HTTPHandler) handleCreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	log.Printf("HTTP Adapter: Received CreateOrder request for customer %s", req.CustomerID)
	order, err := h.orderSvc.CreateOrder(req.CustomerID, req.Items, req.PaymentToken)

	if err != nil {
		log.Printf("HTTP Adapter: CreateOrder failed: %v", err)
		if strings.Contains(err.Error(), "not enough stock") {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (h *HTTPHandler) handleGetOrder(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}
	log.Printf("HTTP Adapter: Received GetOrder request for ID %s", id)
	order, err := h.orderSvc.GetOrderByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *HTTPHandler) handleAddProduct(c *gin.Context) {
	var req AddProductRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("HTTP Adapter: Received AddProduct request for %s", req.Name)
	product, err := h.inventorySvc.AddProduct(req.Name, req.Price, req.Stock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}

func (h *HTTPHandler) handleGetProduct(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id parameter"})
		return
	}
	log.Printf("HTTP Adapter: Received GetProduct request for ID %s", id)
	product, err := h.inventorySvc.GetProduct(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, product)
}

func (h *HTTPHandler) handleGetAllProducts(c *gin.Context) {
	log.Printf("HTTP Adapter: Received GetAllProducts request")
	products, err := h.inventorySvc.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}
