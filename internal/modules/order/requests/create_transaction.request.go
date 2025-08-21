package requests

type CreateTransactionRequest struct {
    ScheduleID    uint    `json:"schedule_id" binding:"required"`
    SeatNumbers   []uint  `json:"seat_numbers" binding:"required,min=1"`
    PaymentMethod string  `json:"payment_method" binding:"required,oneof=credit_card e_wallet"`
}