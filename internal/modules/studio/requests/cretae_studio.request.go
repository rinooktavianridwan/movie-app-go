package requests

type CreateStudioRequest struct {
	Name         string `json:"name" binding:"required"`
	SeatCapacity uint   `json:"capacity" binding:"required"`
	FacilityIDs  []uint `json:"facility_ids"`
}
