package email

import (
	"log"
)

func SendOrderConfirmation(email string, orderID int64) {
	log.Printf("Sending order confirmation email to %s for order ID %d", email, orderID)
	// Simulate email sending
}
