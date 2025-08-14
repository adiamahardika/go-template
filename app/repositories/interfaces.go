package repositories

import "monitoring-service/app/models"

type UserRepositoryInterface interface {
	GetAllUsers(limit, offset int) ([]models.User, int64, error)
	GetUserByID(id int) (*models.User, error)
	EmailExists(email string) (bool, error)
	CreateUser(user models.User) (*models.User, error)
	GetRoleByName(name string) (*models.Role, error)
	AssignRole(userRole models.UserRole) error
}

type ShippingPaymentRepositoryInterface interface {
	CreateShippingMethod(method models.ShippingMethod) (*models.ShippingMethod, error)
	GetShippingMethods(filter models.ShippingMethodFilter) ([]models.ShippingMethod, int64, error)
	GetShippingMethodByID(id int) (*models.ShippingMethod, error)
	UpdateShippingMethod(method models.ShippingMethod) error
	DeleteShippingMethod(id int) error

	CreatePaymentMethod(method models.PaymentMethod) (*models.PaymentMethod, error)
	GetPaymentMethods(filter models.PaymentMethodFilter) ([]models.PaymentMethod, int64, error)
	GetPaymentMethodByID(id int) (*models.PaymentMethod, error)
	UpdatePaymentMethod(method models.PaymentMethod) error
	DeletePaymentMethod(id int) error
}
