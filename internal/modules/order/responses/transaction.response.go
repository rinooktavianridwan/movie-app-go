package responses

import (
    "movie-app-go/internal/models"
    "time"
)

type TransactionResponse struct {
    ID            uint      `json:"id"`
    UserID        uint      `json:"user_id"`
    TotalAmount   float64   `json:"total_amount"`
    PaymentMethod string    `json:"payment_method"`
    PaymentStatus string    `json:"payment_status"`
    CreatedAt     time.Time `json:"created_at"`
    User          UserInfo  `json:"user"`
    Tickets       []TicketInfo `json:"tickets"`
}

type UserInfo struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

type TicketInfo struct {
    ID         uint    `json:"id"`
    SeatNumber uint    `json:"seat_number"`
    Status     string  `json:"status"`
    Price      float64 `json:"price"`
    Schedule   ScheduleInfo `json:"schedule"`
}

type ScheduleInfo struct {
    ID        uint      `json:"id"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
    Date      time.Time `json:"date"`
    Movie     MovieInfo `json:"movie"`
    Studio    StudioInfo `json:"studio"`
}

type MovieInfo struct {
    ID       uint   `json:"id"`
    Title    string `json:"title"`
    Duration uint   `json:"duration"`
}

type StudioInfo struct {
    ID           uint   `json:"id"`
    Name         string `json:"name"`
    SeatCapacity uint   `json:"seat_capacity"`
}

type PaginatedTransactionResponse struct {
    Page      int                   `json:"page"`
    PerPage   int                   `json:"per_page"`
    Total     int64                 `json:"total"`
    TotalPage int                   `json:"total_page"`
    Data      []TransactionResponse `json:"data"`
}

func ToTransactionResponse(transaction *models.Transaction) TransactionResponse {
    tickets := make([]TicketInfo, len(transaction.Tickets))
    for i, ticket := range transaction.Tickets {
        tickets[i] = TicketInfo{
            ID:         ticket.ID,
            SeatNumber: ticket.SeatNumber,
            Status:     ticket.Status,
            Price:      ticket.Price,
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

    return TransactionResponse{
        ID:            transaction.ID,
        UserID:        transaction.UserID,
        TotalAmount:   transaction.TotalAmount,
        PaymentMethod: transaction.PaymentMethod,
        PaymentStatus: transaction.PaymentStatus,
        CreatedAt:     transaction.CreatedAt,
        User: UserInfo{
            ID:    transaction.User.ID,
            Name:  transaction.User.Name,
            Email: transaction.User.Email,
        },
        Tickets: tickets,
    }
}

func ToTransactionResponses(transactions []models.Transaction) []TransactionResponse {
    resp := make([]TransactionResponse, len(transactions))
    for i, t := range transactions {
        resp[i] = ToTransactionResponse(&t)
    }
    return resp
}