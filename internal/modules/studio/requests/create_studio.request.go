package requests

type CreateStudioRequest struct {
	Name         string `json:"name" binding:"required"`
	SeatCapacity uint   `json:"seat_capacity" binding:"required"`
	FacilityIDs  []uint `json:"facility_ids"`
}
