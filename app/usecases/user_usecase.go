package usecases

import (
	"math"
	"monitoring-service/app/models"
)

type userUsecase usecase

type UserUsecaseInterface interface {
	GetAllUsers(request models.GetUsersRequest) ([]models.UserResponse, models.Pagination, error)
	GetUserByID(id int) (*models.UserResponse, error)
}

func (u *userUsecase) GetAllUsers(request models.GetUsersRequest) ([]models.UserResponse, models.Pagination, error) {
	// Set default pagination values
	if request.Page <= 0 {
		request.Page = 1
	}
	if request.PageSize <= 0 {
		request.PageSize = 10
	}

	offset := (request.Page - 1) * request.PageSize

	users, total, err := u.Options.Repository.User.GetAllUsers(request.PageSize, offset)
	if err != nil {
		return nil, models.Pagination{}, err
	}

	// Convert to response DTOs
	userResponses := models.ToUsersResponse(users)

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(request.PageSize)))

	pagination := models.Pagination{
		Page:      request.Page,
		PageSize:  request.PageSize,
		Total:     int(total),
		TotalPage: totalPages,
	}

	return userResponses, pagination, nil
}

func (u *userUsecase) GetUserByID(id int) (*models.UserResponse, error) {
	user, err := u.Options.Repository.User.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Convert to response DTO
	userResponse := user.ToUserResponse()
	return &userResponse, nil
}
