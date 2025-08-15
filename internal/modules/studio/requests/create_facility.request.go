package requests

type CreateFacilityRequest struct {
    Name string `json:"name" binding:"required"`
}