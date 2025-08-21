package requests

type ProcessPaymentRequest struct {
    PaymentStatus string `json:"payment_status" binding:"required,oneof=success failed"`
    PaymentNote   string `json:"payment_note,omitempty"`
}