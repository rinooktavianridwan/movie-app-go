package responses

import (
	"movie-app-go/internal/models"
)

type FacilityResponse struct {
    ID   uint   `json:"id"`
    Name string `json:"name"`
}

type PaginatedFacilityResponse struct {
    Page      int               `json:"page"`
    PerPage   int               `json:"per_page"`
    Total     int64             `json:"total"`
    TotalPage int               `json:"total_page"`
    Data      []FacilityResponse `json:"data"`
}

func ToFacilityResponse(facility *models.Facility) FacilityResponse {
    return FacilityResponse{
        ID:   facility.ID,
        Name: facility.Name,
    }
}

func ToFacilityResponses(facilities []models.Facility) []FacilityResponse {
    resp := make([]FacilityResponse, len(facilities))
    for i, f := range facilities {
        resp[i] = ToFacilityResponse(&f)
    }
    return resp
}