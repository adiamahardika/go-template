package repositories

import (
	"context"
	"monitoring-service/app/models"
)

type userRolesRepository repository

type UserRolesRepositoryInterface interface {
	GetUserRoles(ctx context.Context, userID int) ([]models.Role, error)
	AssignRoleToUser(ctx context.Context, userID int, roleID int) error
	GetRoleByName(ctx context.Context, roleName string) (*models.Role, error)
}

func (r *userRolesRepository) GetUserRoles(ctx context.Context, userID int) ([]models.Role, error) {
	var roles []models.Role

	if err := r.Options.Postgres.
		Model(&models.UserRole{}).
		Select("roles.*").
		Joins("JOIN roles ON roles.id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error; err != nil {
		return nil, err
	}

	return roles, nil
}

func (r *userRolesRepository) AssignRoleToUser(ctx context.Context, userID int, roleID int) error {
	userRole := models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	if err := r.Options.Postgres.Create(&userRole).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRolesRepository) GetRoleByName(ctx context.Context, roleName string) (*models.Role, error) {
	var role models.Role

	if err := r.Options.Postgres.
		Where("name = ? AND deleted_at IS NULL", roleName).
		First(&role).Error; err != nil {
		return nil, err
	}

	return &role, nil
}
