package usecases

import (
	"math"

	"github.com/golang-jwt/jwt"

	"monitoring-service/app/models"
	"monitoring-service/pkg/customerror"
	"monitoring-service/pkg/utils"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase usecase

type UserUsecaseInterface interface {
	GetAllUsers(request models.GetUsersRequest) ([]models.UserResponse, models.Pagination, error)
	GetUserByID(id int) (*models.UserResponse, error)
	Register(request models.RegisterRequest) (string, error)
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

	newUser, err := u.Options.Repository.User.CreateUser(user)
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
