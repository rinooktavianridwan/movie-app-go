package responses

import (
    "movie-app-go/internal/models"
)

type RoleResponse struct {
    ID          uint                    `json:"id"`
    Name        string                  `json:"name"`
    Description string                  `json:"description"`
    Permissions []PermissionResponse    `json:"permissions,omitempty"`
}

type PaginatedRoleResponse struct {
    Page      int            `json:"page"`
    PerPage   int            `json:"per_page"`
    Total     int64          `json:"total"`
    TotalPage int            `json:"total_page"`
    Data      []RoleResponse `json:"data"`
}

func ToRoleResponse(role *models.Role) RoleResponse {
    permissions := make([]PermissionResponse, len(role.Permissions))
    for i, p := range role.Permissions {
        permissions[i] = PermissionResponse{
            ID:          p.ID,
            Name:        p.Name,
            Resource:    p.Resource,
            Action:      p.Action,
            Description: p.Description,
        }
    }

    return RoleResponse{
        ID:          role.ID,
        Name:        role.Name,
        Description: role.Description,
        Permissions: permissions,
    }
}

func ToRoleResponses(roles []models.Role) []RoleResponse {
    resp := make([]RoleResponse, len(roles))
    for i, role := range roles {
        resp[i] = ToRoleResponse(&role)
    }
    return resp
}