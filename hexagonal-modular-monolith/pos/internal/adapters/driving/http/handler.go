package http

import (
	"log"
	"net/http"
	"pos/internal/modules/auth"
	"pos/internal/modules/inventory"
	"pos/internal/modules/orders"
	"strings"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	EmailOrUsername string `json:"email_or_username" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type AuthResponse struct {
	User  *auth.User `json:"user"`
	Token string     `json:"token"`
}

type HTTPHandler struct {
	orderSvc       orders.OrderService
	inventorySvc   inventory.InventoryService
	authSvc        auth.AuthService
	authMiddleware *AuthMiddleware
}

func NewHTTPHandler(os orders.OrderService, is inventory.InventoryService, as auth.AuthService, am *AuthMiddleware) *HTTPHandler {
	return &HTTPHandler{
		orderSvc:       os,
		inventorySvc:   is,
		authSvc:        as,
		authMiddleware: am,
	}
}
func (h *HTTPHandler) RegisterRoutes(gin *gin.Engine) {
	// Auth routes (no auth required)
	gin.POST("/auth/login", h.handleLogin)
	gin.POST("/auth/register", h.handleRegister)

	// Protected routes with RBAC
	auth := h.authMiddleware.RequireAuth()

	// Order routes - require authentication + specific permissions
	gin.POST("/order/create", auth, h.authMiddleware.RequirePermission("order", "create"), h.handleCreateOrder)
	gin.GET("/order/get", auth, h.authMiddleware.RequirePermission("order", "read"), h.handleGetOrder)

	// Product routes - different permissions for different actions
	gin.POST("/product/add", auth, h.authMiddleware.RequirePermission("product", "create"), h.handleAddProduct)
	gin.GET("/product/get", auth, h.authMiddleware.RequirePermission("product", "read"), h.handleGetProduct)
	gin.GET("/product/all", auth, h.authMiddleware.RequirePermission("product", "read"), h.handleGetAllProducts)
}

func (h *HTTPHandler) handleLogin(c *gin.Context) {
	log.Printf("Login: Received login request for: %s", c.Request.URL.Path)
	var req LoginRequest
	if err := c.BindJSON(&req); err != nil {
		log.Printf("Login: Failed to bind JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	log.Printf("Login: Parsed request - email/username: %s", req.EmailOrUsername)

	user, err := h.authSvc.AuthenticateUser(c.Request.Context(), req.EmailOrUsername, req.Password)
	if err != nil {
		log.Printf("Login: Authentication failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := h.authMiddleware.GenerateToken(user)
	if err != nil {
		log.Printf("Login: Token generation failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	log.Printf("Login: Success for user %s", user.Username)
	c.JSON(http.StatusOK, AuthResponse{
		User:  user,
		Token: token,
	})
}

func (h *HTTPHandler) handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Create user through auth service
	user, err := h.authSvc.CreateUser(c.Request.Context(), req.Email, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token
	token, err := h.authMiddleware.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusCreated, AuthResponse{
		User:  user,
		Token: token,
	})
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
