package repositories

import (
	"errors"

	"monitoring-service/app/models"

	"gorm.io/gorm"
)

type userRepository repository

type UserRepositoryInterface interface {
	GetAllUsers(limit, offset int) ([]models.User, int64, error)
	GetUserByID(id int) (*models.User, error)
	EmailExists(email string) (bool, error)
	CreateUser(user models.User) (*models.User, error)
	GetRoleByName(name string) (*models.Role, error)
	AssignRole(userRole models.UserRole) error
	CheckEmailExists(email string) (bool, error)
	GetUserByEmail(email string) (*models.User, error)
	GetActiveCartByUserID(userID int) (*models.Cart, error)
}

func (r *userRepository) GetAllUsers(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	if err := r.Options.Postgres.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.Options.Postgres.
		Preload("UserRoles").
		Preload("UserRoles.Role").
		Limit(limit).
		Offset(offset).
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) GetUserByID(id int) (*models.User, error) {
	var user models.User
	if err := r.Options.Postgres.
		Preload("UserRoles").
		Preload("UserRoles.Role").
		Preload("Carts.CartItems.Product").
		Preload("Orders").
		First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.Options.Postgres.Unscoped().
		Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error
	return count > 0, err
}

func (r *userRepository) CreateUser(user models.User) (*models.User, error) {
	err := r.Options.Postgres.Create(&user).Error
	return &user, err
}

func (r *userRepository) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.Options.Postgres.Where("name = ?", name).First(&role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &role, err
}

func (r *userRepository) AssignRole(userRole models.UserRole) error {
	return r.Options.Postgres.Create(&userRole).Error
}

func (r *userRepository) CheckEmailExists(email string) (bool, error) {
	var count int64
	err := r.Options.Postgres.Model(&models.User{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	if err := r.Options.Postgres.
		Preload("UserRoles").
		Preload("UserRoles.Role").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetActiveCartByUserID(userID int) (*models.Cart, error) {
	var cart models.Cart
	err := r.Options.Postgres.Preload("CartItems.Product").
	Where("user_id = ?", userID).First(&cart).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cart, nil
}
