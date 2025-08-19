package requests

type CreateGenreRequest struct {
	Name string `json:"name" binding:"required"`
}
