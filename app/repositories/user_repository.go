package repositories

import (
	"monitoring-service/app/models"
)

type userRepository repository

type UserRepositoryInterface interface {
	GetAllUsers(limit, offset int) ([]models.User, int64, error)
	GetUserByID(id int) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error) // Tambahkan method baru
}

func (r *userRepository) GetAllUsers(limit, offset int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// Count total users
	if err := r.Options.Postgres.Model(&models.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get users with pagination and preload relationships
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
		Preload("Carts").
		Preload("Orders").
		First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// Tambahkan implementasi baru
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

