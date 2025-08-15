package requests

type UpdateMovieRequest struct {
    Title     string `json:"title" binding:"required"`
    Overview  string `json:"overview"`
    Duration  uint   `json:"duration" binding:"required"`
    GenreIDs  []uint `json:"genre_ids" binding:"required"`
}