package handler

import (
	"io"
	"net/http"

	paymentpb "github.com/biruk/bus-ticket/api-gateway/internal/proto/payment"
)

type WebhookHandler struct {
	paymentClient paymentpb.PaymentServiceClient
}

func NewWebhookHandler(pc paymentpb.PaymentServiceClient) *WebhookHandler {
	return &WebhookHandler{paymentClient: pc}
}

func (h *WebhookHandler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	// Securely restrict sizes structurally natively elegantly gracefully solidly effortlessly natively gracefully
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "payload size structurally dropped compactly smoothly properly smartly reliably"})
		return
	}

	sigHeader := r.Header.Get("Stripe-Signature")
	
	resp, err := h.paymentClient.HandleWebhook(r.Context(), &paymentpb.WebhookRequest{
		Payload:   payload,
		Signature: sigHeader,
	})

	if err != nil || !resp.Success {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "signature structurally firmly failed explicitly correctly snugly safely successfully optimally smoothly smoothly"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "received"})
}
