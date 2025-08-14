package usecases

import (
	"errors"
	"monitoring-service/app/models"
	"monitoring-service/pkg/customerror"
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type UserUsecaseInterface interface {
	GetAllUsers(request models.GetUsersRequest) ([]models.UserResponse, models.Pagination, error)
	GetUserByID(id int) (*models.UserResponse, error)
	Register(request models.RegisterRequest) (string, error)
}

type userUsecase struct {
	*usecase
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

	users, total, err := u.options.Repository.User.GetAllUsers(request.PageSize, offset)
	if err != nil {
		return nil, models.Pagination{}, err
	}

	responses := models.ToUsersResponse(users)

	totalPages := 0
	if request.PageSize > 0 {
		totalPages = int((total + int64(request.PageSize) - 1) / int64(request.PageSize))
	}

	pagination := models.Pagination{
		Page:      request.Page,
		PageSize:  request.PageSize,
		Total:     int(total),
		TotalPage: totalPages,
	}

	return responses, pagination, nil
}

func (u *userUsecase) GetUserByID(id int) (*models.UserResponse, error) {
	user, err := u.options.Repository.User.GetUserByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, customerror.NewNotFoundError("user not found")
		}
		return nil, err
	}

	response := user.ToUserResponse()
	return &response, nil
}

func (u *userUsecase) Register(request models.RegisterRequest) (string, error) {
	// Check email uniqueness
	exists, err := u.options.Repository.User.EmailExists(request.Email)
	if err != nil {
		return "", err
	}
	if exists {
		return "", customerror.NewConflictError("email already exists")
	}

	// Create user
	user := models.User{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password, // Note: Password should be hashed before saving
	}

	newUser, err := u.options.Repository.User.CreateUser(user)
	if err != nil {
		return "", err
	}

	// Assign shopper role
	shopperRole, err := u.options.Repository.User.GetRoleByName("shopper")
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

	if err := u.options.Repository.User.AssignRole(userRole); err != nil {
		return "", err
	}

	// Generate JWT token
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = newUser.ID
	claims["role"] = "shopper"
	claims["exp"] = time.Now().Add(time.Hour * 24 * time.Duration(u.options.Config.JWTExpireTime)).Unix()

	tokenString, err := token.SignedString([]byte(u.options.Config.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
