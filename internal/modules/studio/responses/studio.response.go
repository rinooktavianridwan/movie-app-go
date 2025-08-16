package responses

import (
	"movie-app-go/internal/models"
)

type StudioResponse struct {
    ID           uint   `json:"id"`
    Name         string `json:"name"`
    SeatCapacity uint   `json:"seat_capacity"`
    Facilities   []FacilityResponse `json:"facilities"`
}

type PaginatedStudioResponse struct {
    Page      int              `json:"page"`
    PerPage   int              `json:"per_page"`
    Total     int64            `json:"total"`
    TotalPage int              `json:"total_page"`
    Data      []StudioResponse `json:"data"`
}

func ToStudioResponse(studio *models.Studio) StudioResponse {
    facilities := make([]FacilityResponse, len(studio.FacilityStudios))
    for i, fs := range studio.FacilityStudios {
        facilities[i] = FacilityResponse{
            ID:   fs.Facility.ID,
            Name: fs.Facility.Name,
        }
    }
    return StudioResponse{
        ID:           studio.ID,
        Name:         studio.Name,
        SeatCapacity: studio.SeatCapacity,
        Facilities:   facilities,
    }
}

func ToStudioResponses(studios []models.Studio) []StudioResponse {
    resp := make([]StudioResponse, len(studios))
    for i, s := range studios {
        resp[i] = ToStudioResponse(&s)
    }
    return resp
}