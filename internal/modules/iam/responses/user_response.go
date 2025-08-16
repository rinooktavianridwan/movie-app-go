package responses

import (
	"movie-app-go/internal/models"
)

type UserResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	IsAdmin bool   `json:"is_admin"`
}

type PaginatedUserResponse struct {
	Page      int            `json:"page"`
	PerPage   int            `json:"per_page"`
	Total     int64          `json:"total"`
	TotalPage int            `json:"total_page"`
	Data      []UserResponse `json:"data"`
}

func ToUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:      user.ID,
		Name:    user.Name,
		Email:   user.Email,
		IsAdmin: user.IsAdmin,
	}
}

func ToUserResponses(users []models.User) []UserResponse {
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ToUserResponse(&user)
	}
	return userResponses
}
