package user

import (
	"food-delivery-workshop/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *models.User) error
	FindByEmail(email string, user *models.User) error
	FindByID(userID uint, user *models.User) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) FindByEmail(email string, user *models.User) error {
	if err := r.db.Where("email = ?", email).First(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) FindByID(userID uint, user *models.User) error {
	if err := r.db.Where("id = ?", userID).First(user).Error; err != nil {
		return err
	}
	return nil
}
