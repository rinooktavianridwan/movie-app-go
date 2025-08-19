package responses

import (
	"movie-app-go/internal/models"
)

type GenreResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type PaginatedGenreResponse struct {
	Page      int             `json:"page"`
	PerPage   int             `json:"per_page"`
	Total     int64           `json:"total"`
	TotalPage int             `json:"total_page"`
	Data      []GenreResponse `json:"data"`
}

func ToGenreResponse(genre *models.Genre) GenreResponse {
	return GenreResponse{
		ID:   genre.ID,
		Name: genre.Name,
	}
}

func ToGenreResponses(genres []models.Genre) []GenreResponse {
	resp := make([]GenreResponse, len(genres))
	for i, g := range genres {
		resp[i] = ToGenreResponse(&g)
	}
	return resp
}
