package responses

import (
    "movie-app-go/internal/models"
    "time"
)

type TicketResponse struct {
    ID            uint      `json:"id"`
    TransactionID uint      `json:"transaction_id"`
    ScheduleID    uint      `json:"schedule_id"`
    SeatNumber    uint      `json:"seat_number"`
    Status        string    `json:"status"`
    Price         float64   `json:"price"`
    CreatedAt     time.Time `json:"created_at"`
    Transaction   TransactionInfo `json:"transaction"`
    Schedule      ScheduleInfo    `json:"schedule"`
}

type TransactionInfo struct {
    ID            uint   `json:"id"`
    PaymentMethod string `json:"payment_method"`
    PaymentStatus string `json:"payment_status"`
    User          UserInfo `json:"user"`
}

type PaginatedTicketResponse struct {
    Page      int              `json:"page"`
    PerPage   int              `json:"per_page"`
    Total     int64            `json:"total"`
    TotalPage int              `json:"total_page"`
    Data      []TicketResponse `json:"data"`
}

func ToTicketResponse(ticket *models.Ticket) TicketResponse {
    return TicketResponse{
        ID:            ticket.ID,
        TransactionID: ticket.TransactionID,
        ScheduleID:    ticket.ScheduleID,
        SeatNumber:    ticket.SeatNumber,
        Status:        ticket.Status,
        Price:         ticket.Price,
        CreatedAt:     ticket.CreatedAt,
        Transaction: TransactionInfo{
            ID:            ticket.Transaction.ID,
            PaymentMethod: ticket.Transaction.PaymentMethod,
            PaymentStatus: ticket.Transaction.PaymentStatus,
            User: UserInfo{
                ID:    ticket.Transaction.User.ID,
                Name:  ticket.Transaction.User.Name,
                Email: ticket.Transaction.User.Email,
            },
        },
        Schedule: ScheduleInfo{
            ID:        ticket.Schedule.ID,
            StartTime: ticket.Schedule.StartTime,
            EndTime:   ticket.Schedule.EndTime,
            Date:      ticket.Schedule.Date,
            Movie: MovieInfo{
                ID:       ticket.Schedule.Movie.ID,
                Title:    ticket.Schedule.Movie.Title,
                Duration: ticket.Schedule.Movie.Duration,
            },
            Studio: StudioInfo{
                ID:           ticket.Schedule.Studio.ID,
                Name:         ticket.Schedule.Studio.Name,
                SeatCapacity: ticket.Schedule.Studio.SeatCapacity,
            },
        },
    }
}

func ToTicketResponses(tickets []models.Ticket) []TicketResponse {
    resp := make([]TicketResponse, len(tickets))
    for i, t := range tickets {
        resp[i] = ToTicketResponse(&t)
    }
    return resp
}