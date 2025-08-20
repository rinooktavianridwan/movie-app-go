package responses

import (
    "movie-app-go/internal/models"
    "time"
)

type ScheduleResponse struct {
    ID        uint      `json:"id"`
    MovieID   uint      `json:"movie_id"`
    StudioID  uint      `json:"studio_id"`
    StartTime time.Time `json:"start_time"`
    EndTime   time.Time `json:"end_time"`
    Date      time.Time `json:"date"`
    Price     float64   `json:"price"`
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

type PaginatedScheduleResponse struct {
    Page      int                `json:"page"`
    PerPage   int                `json:"per_page"`
    Total     int64              `json:"total"`
    TotalPage int                `json:"total_page"`
    Data      []ScheduleResponse `json:"data"`
}

func ToScheduleResponse(schedule *models.Schedule) ScheduleResponse {
    return ScheduleResponse{
        ID:        schedule.ID,
        MovieID:   schedule.MovieID,
        StudioID:  schedule.StudioID,
        StartTime: schedule.StartTime,
        EndTime:   schedule.EndTime,
        Date:      schedule.Date,
        Price:     schedule.Price,
        Movie: MovieInfo{
            ID:       schedule.Movie.ID,
            Title:    schedule.Movie.Title,
            Duration: schedule.Movie.Duration,
        },
        Studio: StudioInfo{
            ID:           schedule.Studio.ID,
            Name:         schedule.Studio.Name,
            SeatCapacity: schedule.Studio.SeatCapacity,
        },
    }
}

func ToScheduleResponses(schedules []models.Schedule) []ScheduleResponse {
    resp := make([]ScheduleResponse, len(schedules))
    for i, s := range schedules {
        resp[i] = ToScheduleResponse(&s)
    }
    return resp
}