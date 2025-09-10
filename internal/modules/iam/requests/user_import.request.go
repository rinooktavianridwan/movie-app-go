package requests

type ImportUserRequest struct {
    SheetName string `form:"sheet_name" binding:"omitempty"`
}
