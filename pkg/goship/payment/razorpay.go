package payment

import (
	"github.com/razorpay/razorpay-go"
)

// RazorpayClient wraps the Razorpay API client.
type RazorpayClient struct {
	Client *razorpay.Client
}

// NewRazorpayClient initializes a Razorpay client.
//
// Example:
//
//	client := NewRazorpayClient("rzp_test_xxx", "secret")
func NewRazorpayClient(key, secret string) *RazorpayClient {
	return &RazorpayClient{
		Client: razorpay.NewClient(key, secret),
	}
}

// CreateOrder creates a new Razorpay order.
//
// Example:
//
//	order, err := client.CreateOrder(500, "INR", "receipt#1")
func (r *RazorpayClient) CreateOrder(amount int, currency, receipt string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"amount":   amount * 100, // Razorpay expects amount in paise
		"currency": currency,
		"receipt":  receipt,
	}
	return r.Client.Order.Create(data, nil)
}
