package payment

import (
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/charge"
)

// InitStripe initializes the Stripe client with the secret key.
//
// Example:
//
//	InitStripe("sk_test_...")
func InitStripe(secretKey string) {
	stripe.Key = secretKey
}

// Charge creates a Stripe charge.
//
// Parameters:
//   - amount: Amount in smallest currency unit (e.g., 500 = $5.00).
//   - currency: Currency code (e.g., "usd").
//   - source: Payment source (e.g., token from client).
//   - description: Charge description.
//
// Example:
//
//	err := Charge(500, "usd", "tok_visa", "Test payment")
func Charge(amount int64, currency, source, description string) error {
	_, err := charge.New(&stripe.ChargeParams{
		Amount:      stripe.Int64(amount),
		Currency:    stripe.String(currency),
		Source:      &stripe.PaymentSourceSourceParams{Token: stripe.String(source)},
		Description: stripe.String(description),
	})
	return err
}
