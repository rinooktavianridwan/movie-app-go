package requests

type UpdateFacilityRequest struct {
    Name string `json:"name" binding:"required"`
}