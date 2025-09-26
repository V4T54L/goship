package payment

import (
	"github.com/braintree-go/braintree-go"
)

// NewBraintreeClient initializes the Braintree client.
//
// Parameters:
//   - env: "sandbox" or "production"
//   - merchantID, publicKey, privateKey: credentials
//
// Example:
//
//	client := NewBraintreeClient("sandbox", "merchant", "pub", "priv")
func NewBraintreeClient(env, merchantID, publicKey, privateKey string) *braintree.Braintree {
	var btEnv braintree.Environment
	if env == "production" {
		btEnv = braintree.Production
	} else {
		btEnv = braintree.Sandbox
	}
	return braintree.New(btEnv, merchantID, publicKey, privateKey)
}
