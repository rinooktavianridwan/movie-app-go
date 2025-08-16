package requests

type RegisterRequest struct {
    Name     string `json:"name" binding:"required"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required"`
    IsAdmin  *bool  `json:"is_admin,omitempty"`
}