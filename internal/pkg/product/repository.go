package product

import (
	"food-delivery-workshop/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(product *models.Product) error
	Update(product *models.Product) error
	FindByID(id uint, product *models.Product) error
	FindAll() ([]models.Product, error)
	Delete(id uint) error
	Preload(product interface{}) error
	FindByProductName(name string) (*models.Product, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r* repository) Preload(product interface{}) error {
	return r.db.Preload("Promotion").
	Find(product).Error
}

func (r *repository) Create(product *models.Product) error {
	if err := r.db.Create(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) FindByID(id uint, product *models.Product) error {
	if err := r.db.Where("id = ?", id).First(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Where("deleted_at IS NULL").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repository) Update(product *models.Product) error {
	if err := r.db.Save(product).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(id uint) error {
	if err := r.db.Where("id =?", id).Delete(&models.Product{}).Error; err!= nil {
        return err
    }

	return nil
}

func (r *repository) FindByProductName(name string) (*models.Product, error) {
	var product models.Product
    if err := r.db.Where("name =?", name).First(&product).Error; err != nil {
        return nil, err
    }
    return &product, nil
}
