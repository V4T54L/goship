package payment

import (
	"github.com/plutov/paypal/v4"
)

// NewPayPalClient initializes the PayPal client.
//
// Parameters:
//   - clientID: PayPal client ID.
//   - secret: PayPal secret.
//   - isLive: true for live, false for sandbox.
//
// Example:
//
//	client, err := NewPayPalClient("id", "secret", false)
func NewPayPalClient(clientID, secret string, isLive bool) (*paypal.Client, error) {
	env := paypal.APIBaseSandBox
	if isLive {
		env = paypal.APIBaseLive
	}
	return paypal.NewClient(clientID, secret, env)
}
