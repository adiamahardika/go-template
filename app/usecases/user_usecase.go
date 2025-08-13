package usecases

import (
	"math"
	"monitoring-service/app/models"
	"monitoring-service/pkg/utils"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase usecase

type UserUsecaseInterface interface {
	GetAllUsers(request models.GetUsersRequest) ([]models.UserResponse, models.Pagination, error)
	GetUserByID(id int) (*models.UserResponse, error)
	Login(request models.LoginRequest) (*models.LoginResponse, error) // Tambahkan method baru
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

// Tambahkan implementasi baru
func (u *userUsecase) Login(request models.LoginRequest) (*models.LoginResponse, error) {
	// Cari user berdasarkan email
	user, err := u.Options.Repository.User.GetUserByEmail(request.Email)
	if err != nil {
		return nil, errors.Wrap(err, "invalid email or password")
	}

	// Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Dapatkan role user (asumsi ada relasi UserRoles)
	var role string
	if len(user.UserRoles) > 0 {
		role = user.UserRoles[0].Role.Name
	} else {
		role = "shopper" // Default role
	}

	// Buat token JWT
	expireTime := time.Now().Add(time.Hour * time.Duration(u.Options.Config.JWTExpireTime))
	token, err := utils.GenerateJWTToken(user.ID, role, u.Options.Config.JWTSecret, expireTime)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate token")
	}

	// Siapkan response
	response := &models.LoginResponse{
		Token:     token,
		ExpiresAt: expireTime,
		User: models.UserAuth{
			ID:    user.ID,
			Email: user.Email,
			Role:  role,
		},
	}

	return response, nil
}
