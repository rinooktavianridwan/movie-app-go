package requests

import "time"

type UpdateScheduleRequest struct {
    MovieID   *uint      `json:"movie_id,omitempty"`
    StudioID  *uint      `json:"studio_id,omitempty"`  
    StartTime *time.Time `json:"start_time,omitempty"`
    Date      *time.Time `json:"date,omitempty"`
    Price     *float64   `json:"price,omitempty" binding:"min=0"`
}
