package requests

type UpdateTicketRequest struct {
    Status string `json:"status" binding:"required,oneof=active used cancelled"`
}