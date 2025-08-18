package requests

import "time"

type UpdateScheduleRequest struct {
    MovieID   uint      `json:"movie_id" binding:"required"`
    StudioID  uint      `json:"studio_id" binding:"required"`
    StartTime time.Time `json:"start_time" binding:"required"`
    Date      time.Time `json:"date" binding:"required"`
    Price     float64   `json:"price" binding:"required,min=0"`
}