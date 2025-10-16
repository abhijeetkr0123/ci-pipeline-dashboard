package handlers

import (
	"fmt"
	"net/http"
)

// WebhookHandler handles POST requests
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Webhook received!")
}
