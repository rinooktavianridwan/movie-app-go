package responses

import (
	"movie-app-go/internal/models"
)

type MovieResponse struct {
	ID       uint            `json:"id"`
	Title    string          `json:"title"`
	Overview string          `json:"overview"`
	Duration uint            `json:"duration"`
	Genres   []GenreResponse `json:"genres"`
}

type PaginatedMovieResponse struct {
	Page      int             `json:"page"`
	PerPage   int             `json:"per_page"`
	Total     int64           `json:"total"`
	TotalPage int             `json:"total_page"`
	Data      []MovieResponse `json:"data"`
}

// Converter functions
func ToMovieResponse(movie *models.Movie) MovieResponse {
	genres := make([]GenreResponse, len(movie.MovieGenres))
	for i, mg := range movie.MovieGenres {
		genres[i] = GenreResponse{
			ID:   mg.Genre.ID,
			Name: mg.Genre.Name,
		}
	}
	return MovieResponse{
		ID:       movie.ID,
		Title:    movie.Title,
		Overview: movie.Overview,
		Duration: movie.Duration,
		Genres:   genres,
	}
}

func ToMovieResponses(movies []models.Movie) []MovieResponse {
	resp := make([]MovieResponse, len(movies))
	for i, m := range movies {
		resp[i] = ToMovieResponse(&m)
	}
	return resp
}
