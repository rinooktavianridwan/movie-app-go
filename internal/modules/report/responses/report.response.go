package responses

type DailySalesReport struct {
	Date         string  `json:"date"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalTickets int64   `json:"total_tickets"`
}

type MonthlySalesReport struct {
	Month        string  `json:"month"`
	Year         int     `json:"year"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalTickets int64   `json:"total_tickets"`
}

type DailyMovieReport struct {
	Date         string  `json:"date"`
	MovieID      uint    `json:"movie_id"`
	MovieTitle   string  `json:"movie_title"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalTickets int64   `json:"total_tickets"`
}

type DailyStudioReport struct {
	Date         string  `json:"date"`
	StudioID     uint    `json:"studio_id"`
	StudioName   string  `json:"studio_name"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalTickets int64   `json:"total_tickets"`
}

type MonthlyMovieReport struct {
	Month        string  `json:"month"`
	Year         int     `json:"year"`
	MovieID      uint    `json:"movie_id"`
	MovieTitle   string  `json:"movie_title"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalTickets int64   `json:"total_tickets"`
}

type MonthlyStudioReport struct {
	Month        string  `json:"month"`
	Year         int     `json:"year"`
	StudioID     uint    `json:"studio_id"`
	StudioName   string  `json:"studio_name"`
	TotalRevenue float64 `json:"total_revenue"`
	TotalTickets int64   `json:"total_tickets"`
}

type PaginatedDailySalesResponse struct {
    Page      int                `json:"page"`
    PerPage   int                `json:"per_page"`
    Total     int64              `json:"total"`
    TotalPage int                `json:"total_page"`
    Data      []DailySalesReport `json:"data"`
}

type PaginatedMonthlySalesResponse struct {
    Page      int                  `json:"page"`
    PerPage   int                  `json:"per_page"`
    Total     int64                `json:"total"`
    TotalPage int                  `json:"total_page"`
    Data      []MonthlySalesReport `json:"data"`
}

type PaginatedDailyMovieResponse struct {
    Page      int                `json:"page"`
    PerPage   int                `json:"per_page"`
    Total     int64              `json:"total"`
    TotalPage int                `json:"total_page"`
    Data      []DailyMovieReport `json:"data"`
}

type PaginatedDailyStudioResponse struct {
    Page      int                 `json:"page"`
    PerPage   int                 `json:"per_page"`
    Total     int64               `json:"total"`
    TotalPage int                 `json:"total_page"`
    Data      []DailyStudioReport `json:"data"`
}

type PaginatedMonthlyMovieResponse struct {
    Page      int                  `json:"page"`
    PerPage   int                  `json:"per_page"`
    Total     int64                `json:"total"`
    TotalPage int                  `json:"total_page"`
    Data      []MonthlyMovieReport `json:"data"`
}

type PaginatedMonthlyStudioResponse struct {
    Page      int                   `json:"page"`
    PerPage   int                   `json:"per_page"`
    Total     int64                 `json:"total"`
    TotalPage int                   `json:"total_page"`
    Data      []MonthlyStudioReport `json:"data"`
}