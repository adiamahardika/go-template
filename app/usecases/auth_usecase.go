package usecases

import (
	"context"
	"errors"
	"monitoring-service/app/models"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type authUsecase usecase

type AuthUsecaseInterface interface {
	Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error)
	Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error)
}

func (u *authUsecase) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {

	user, err := u.Options.Repository.User.GetUserByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("Invalid email")
		}
		return nil, err
	}

	err = CheckPassword(user.Password, req.Password)
	if err != nil {
		return nil, errors.New("Invalid password")
	}

	roles, err := u.Options.Repository.UserRole.GetUserRoles(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	roleNames := make([]string, len(roles))
	for i, role := range roles {
		roleNames[i] = role.Name
	}

	token, err := u.GenerateToken(*user, roleNames)
	if err != nil {
		return nil, err
	}

	userProfile := models.UserProfile{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Roles: roles,
	}

	expireTime := u.Options.Config.JWT.ExpireTime
	if expireTime <= 0 {
		expireTime = 24
	}

	return &models.AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expireTime * 3600,
		User:        userProfile,
	}, nil
}

func CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *authUsecase) GenerateToken(user models.User, roles []string) (string, error) {
	if u.Options.Config.JWT.Secret == "" {
		return "", errors.New("JWT secret is not set")
	}

	expireTime := time.Duration(u.Options.Config.JWT.ExpireTime) * time.Hour
	if u.Options.Config.JWT.ExpireTime <= 0 {
		expireTime = 24 * time.Hour
	}

	claims := &models.Claims{
		UserID: user.ID,
		Email:  user.Email,
		Roles:  roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   string(rune(user.ID)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.Options.Config.JWT.Secret))
}

func (u *authUsecase) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	// Validate input
	if err := u.validateRegistrationInput(req); err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := u.Options.Repository.User.CheckEmailExists(req.Email)
	if err != nil {
		return nil, errors.New("failed to check email availability")
	}
	if exists {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	createdUser, err := u.Options.Repository.User.CreateUser(*user)
	if err != nil {
		return nil, errors.New("failed to create user")
	}
	user = createdUser

	// Get shopper role
	shopperRole, err := u.Options.Repository.UserRole.GetRoleByName(ctx, "shopper")
	if err != nil {
		return nil, errors.New("shopper role not found")
	}

	// Assign shopper role to user
	if err := u.Options.Repository.UserRole.AssignRoleToUser(ctx, user.ID, shopperRole.ID); err != nil {
		return nil, errors.New("failed to assign role to user")
	}

	// Generate token for immediate login
	roles := []string{"shopper"}
	token, err := u.GenerateToken(*user, roles)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	userProfile := models.UserProfile{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Roles: []models.Role{*shopperRole},
	}

	expireTime := u.Options.Config.JWT.ExpireTime
	if expireTime <= 0 {
		expireTime = 24
	}

	return &models.AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   expireTime * 3600,
		User:        userProfile,
	}, nil
}

func (u *authUsecase) validateRegistrationInput(req models.RegisterRequest) error {
	// Validate name
	if strings.TrimSpace(req.Name) == "" {
		return errors.New("name is required")
	}

	// Validate email
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invalid email format")
	}

	// Validate password
	if err := u.validatePassword(req.Password); err != nil {
		return err
	}

	return nil
}

func (u *authUsecase) validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Check for at least one letter
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	if !hasLetter {
		return errors.New("password must contain at least one letter")
	}

	// Check for at least one number
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}

	return nil
}
