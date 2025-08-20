package dto

import (
	"time"
)

type User struct {
	ID        int
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	UserRoles []UserRole
}

type UserRole struct {
	UserID    int
	RoleID    int
	CreatedAt time.Time
	Role      Role
}

type Role struct {
	ID          int
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}

type UserResponse struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// Relationships
	UserRoles []UserRoleResponse `json:"user_roles,omitempty"`
}

type UserRoleResponse struct {
	UserID    int           `json:"user_id"`
	RoleID    int           `json:"role_id"`
	CreatedAt time.Time     `json:"created_at"`
	Role      *RoleResponse `json:"role,omitempty"`
}

type RoleResponse struct {
	ID          int        `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

func (u *User) ToUserResponse() UserResponse {
	response := UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}


	for _, ur := range u.UserRoles {
		userRole := UserRoleResponse{
			UserID:    ur.UserID,
			RoleID:    ur.RoleID,
			CreatedAt: ur.CreatedAt,
		}


		if ur.Role.ID != 0 {
			userRole.Role = &RoleResponse{
				ID:          ur.Role.ID,
				Name:        ur.Role.Name,
				Description: ur.Role.Description,
				CreatedAt:   ur.Role.CreatedAt,
				UpdatedAt:   ur.Role.UpdatedAt,
				DeletedAt:   ur.Role.DeletedAt,
			}
		}

		response.UserRoles = append(response.UserRoles, userRole)
	}

	return response
}

func ToUsersResponse(users []User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToUserResponse()
	}
	return responses
}


type GetUsersRequest struct {
	Page     int `json:"page" query:"page"`
	PageSize int `json:"page_size" query:"page_size"`
}

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	RoleIDs  []int  `json:"role_ids,omitempty"`
}

type UpdateUserRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	RoleIDs []int  `json:"role_ids,omitempty"`
}
