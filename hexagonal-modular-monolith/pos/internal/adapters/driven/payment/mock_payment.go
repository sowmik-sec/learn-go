package payment

import (
	"errors"
	"log"
	"pos/internal/modules/orders"

	"github.com/google/uuid"
)

type MockPaymentGateway struct{}

func NewMockPaymentGateway() orders.PaymentGateway {
	return &MockPaymentGateway{}
}

func (m *MockPaymentGateway) ProcessPayment(amount float64, token string) (string, error) {
	log.Printf("Payment Adapter: Processing payment of %.2f with token '%s", amount, token)
	if token == "fail-payment" {
		log.Println("Payment Adapter: Simulating payment failure")
		return "", errors.New("invalid card token")
	}
	txID := "txn_" + uuid.NewString()
	log.Printf("Payment Adapter: Payment successful (Tx: %s)", txID)
	return txID, nil
}
