package dto

import (
	"time"
)

// If the User struct is not defined elsewhere, define it here:
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

// UserResponse represents the user data for API responses (without sensitive fields)
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

// ToUserResponse converts User model to UserResponse (excluding password)
func (u *User) ToUserResponse() UserResponse {
	response := UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: u.DeletedAt,
	}

	// Convert user roles
	for _, ur := range u.UserRoles {
		userRole := UserRoleResponse{
			UserID:    ur.UserID,
			RoleID:    ur.RoleID,
			CreatedAt: ur.CreatedAt,
		}

		// Convert role if exists
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

// ToUsersResponse converts slice of User models to slice of UserResponse
func ToUsersResponse(users []User) []UserResponse {
	responses := make([]UserResponse, len(users))
	for i, user := range users {
		responses[i] = user.ToUserResponse()
	}
	return responses
}

// Request DTOs

// GetUsersRequest represents the request for getting users with filters and pagination
type GetUsersRequest struct {
	Page     int `json:"page" query:"page"`
	PageSize int `json:"page_size" query:"page_size"`
}

// CreateUserRequest represents the request for creating a new user
type CreateUserRequest struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	RoleIDs  []int  `json:"role_ids,omitempty"`
}

// UpdateUserRequest represents the request for updating a user
type UpdateUserRequest struct {
	Name    string `json:"name" validate:"required"`
	Email   string `json:"email" validate:"required,email"`
	RoleIDs []int  `json:"role_ids,omitempty"`
}
