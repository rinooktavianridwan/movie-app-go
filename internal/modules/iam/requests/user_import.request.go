package requests

type ImportUserRequest struct {
	SheetName string `json:"sheet_name" binding:"omitempty"`
}
