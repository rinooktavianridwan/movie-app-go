package requests

type UserUpdateRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password,omitempty"`
    IsAdmin  *bool  `json:"is_admin,omitempty"`
}