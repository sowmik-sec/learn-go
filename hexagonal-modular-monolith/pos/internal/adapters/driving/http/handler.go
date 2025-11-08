package http

import (
	"encoding/json"
	"log"
	"net/http"
	"pos/internal/modules/inventory"
	"pos/internal/modules/orders"
	"strings"
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

func (h *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/order/create", h.handleCreateOrder)
	mux.HandleFunc("/order/get", h.handleGetOrder)
	mux.HandleFunc("/product/add", h.handleAddProduct)
	mux.HandleFunc("/product/get", h.handleGetProduct)
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

func (h *HTTPHandler) handleCreateOrder(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("HTTP Adapter: Received CreateOrder request for customer %s", req.CustomerID)
	order, err := h.orderSvc.CreateOrder(req.CustomerID, req.Items, req.PaymentToken)

	if err != nil {
		log.Printf("HTTP Adapter: CreateOrder failed: %v", err)
		if strings.Contains(err.Error(), "not enough stock") {
			httpError(w, err.Error(), http.StatusConflict)
			return
		}
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, order, http.StatusCreated)
}

func (h *HTTPHandler) handleGetOrder(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		httpError(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	log.Printf("HTTP Adapter: Received GetOrder request for ID %s", id)
	order, err := h.orderSvc.GetOrderByID(id)
	if err != nil {
		httpError(w, err.Error(), http.StatusNotFound)
		return
	}
	jsonResponse(w, order, http.StatusOK)
}

func (h *HTTPHandler) handleAddProduct(w http.ResponseWriter, r *http.Request) {
	var req AddProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httpError(w, "invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("HTTP Adapter: Received AddProduct request for %s", req.Name)
	product, err := h.inventorySvc.AddProduct(req.Name, req.Price, req.Stock)
	if err != nil {
		httpError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonResponse(w, product, http.StatusCreated)
}

func (h *HTTPHandler) handleGetProduct(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		httpError(w, "missing id parameter", http.StatusBadRequest)
		return
	}
	log.Printf("HTTP Adapter: Received GetProduct request for ID %s", id)
	product, err := h.inventorySvc.GetProduct(id)
	if err != nil {
		httpError(w, err.Error(), http.StatusNotFound)
		return
	}
	jsonResponse(w, product, http.StatusOK)
}

func httpError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func jsonResponse(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}
