package usecases

import (
	"errors"
	"math"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"monitoring-service/app/models"
	"monitoring-service/pkg/customerror"
)

type userUsecase usecase

type UserUsecaseInterface interface {
	GetAllUsers(request models.GetUsersRequest) ([]models.UserResponse, models.Pagination, error)
	GetUserByID(id int) (*models.UserResponse, error)
	Register(request models.RegisterRequest) (string, error)
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

func (u *userUsecase) Register(request models.RegisterRequest) (string, error) {
	// Check email uniqueness
	exists, err := u.Options.Repository.User.EmailExists(request.Email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", customerror.NewConflictError("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	// Create user
	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: string(hashedPassword),
	}

	newUser, err := u.Options.Repository.User.CreateUser(&user)
	if err != nil {
		return "", err
	}

	// Assign shopper role
	shopperRole, err := u.Options.Repository.User.GetRoleByName("shopper")
	if err != nil {
		return "", err
	}
	if shopperRole == nil {
		return "", errors.New("shopper role not found in database")
	}

	userRole := models.UserRole{
		UserID: newUser.ID,
		RoleID: shopperRole.ID,
	}

	if err := u.Options.Repository.User.AssignRole(userRole); err != nil {
		return "", err
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = newUser.ID
	claims["role"] = "shopper"
	claims["exp"] = time.Now().Add(time.Hour * 24 * time.Duration(u.Options.Config.JWTExpireTime)).Unix()

	tokenString, err := token.SignedString([]byte(u.Options.Config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
