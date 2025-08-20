package requests

type UpdateTransactionRequest struct {
    PaymentStatus string `json:"payment_status" binding:"required,oneof=success failed"`
}