package requests

type UpdateMovieRequest struct {
    Title    string `form:"title" json:"title" binding:"required"`
    Overview string `form:"overview" json:"overview"`
    Duration uint   `form:"duration" json:"duration" binding:"required"`
    GenreIDs []uint `form:"genre_ids" json:"genre_ids" binding:"required"`
}
