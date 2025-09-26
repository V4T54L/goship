package payment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// CreateAdyenPayment creates a payment via Adyen API.
//
// Parameters:
//   - apiKey: Your Adyen API key.
//   - merchantAccount: Your Adyen merchant account.
//   - amount: Amount in minor units (e.g., cents).
//   - currency: Currency code (e.g., "EUR").
//   - reference: Payment reference.
//
// Example:
//
//	err := CreateAdyenPayment("api_key", "merchant123", 500, "EUR", "ref#1")
func CreateAdyenPayment(apiKey, merchantAccount string, amount int, currency, reference string) error {
	payload := map[string]interface{}{
		"amount": map[string]interface{}{
			"value":    amount,
			"currency": currency,
		},
		"merchantAccount": merchantAccount,
		"reference":       reference,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://checkout-test.adyen.com/v68/payments", bytes.NewBuffer(body))
	req.Header.Set("X-API-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("adyen error: %s", resp.Status)
	}
	return nil
}
