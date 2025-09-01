package responses

import (
	"movie-app-go/internal/models"
	"movie-app-go/internal/utils"
	"os"
)

type UserResponse struct {
    ID        uint    `json:"id"`
    Name      string  `json:"name"`
    Email     string  `json:"email"`
    Avatar    *string `json:"avatar"`
    IsAdmin   bool    `json:"is_admin"`
    CreatedAt string  `json:"created_at"`
    UpdatedAt string  `json:"updated_at"`
}

type PaginatedUserResponse struct {
	Page      int            `json:"page"`
	PerPage   int            `json:"per_page"`
	Total     int64          `json:"total"`
	TotalPage int            `json:"total_page"`
	Data      []UserResponse `json:"data"`
}

func ToUserResponse(user *models.User) UserResponse {
	var avatarURL *string
	if user.Avatar != nil && *user.Avatar != "" {
		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "http://localhost:3000"
		}
		fullURL := utils.GetFileURL(*user.Avatar, baseURL)
		avatarURL = &fullURL
	}

	return UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Avatar:    avatarURL,
		IsAdmin:   user.IsAdmin,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func ToUserResponses(users []models.User) []UserResponse {
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = ToUserResponse(&user)
	}
	return userResponses
}
