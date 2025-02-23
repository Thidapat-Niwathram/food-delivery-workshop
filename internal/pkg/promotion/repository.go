package promotion

import (
	"food-delivery-workshop/internal/models"
	"gorm.io/gorm"
)

type Repository interface {
	Create(promotion *models.Promotion) error
	Preload(promotion interface{}) error
	FindByID(id uint, promotion *models.Promotion) error
	FindAll() ([]*models.Promotion, error)
	Update(promotion *models.Promotion) error
	Delete(id uint) error
	FindPromotionByProductID(productID uint) (*models.Promotion, error)
	FindPromotionByCode(code string) (*models.Promotion, error)
	FindPromotionByID(promotionID uint) (*models.Promotion, error)
	DeletePromotionID(promotionID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(promotion *models.Promotion) error {
	if err := r.db.Create(promotion).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) Preload(promotions interface{}) error {
	return r.db.Preload("Product").
		Find(promotions).Error
}

func (r *repository) FindByID(id uint, promotion *models.Promotion) error {
	if err := r.db.Where("id =?", id).Preload("Product").First(promotion).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) FindAll() ([]*models.Promotion, error) {
	var promotions []*models.Promotion
	if err := r.db.Where("deleted_at IS NULL").Preload("Product").Find(&promotions).Error; err != nil {
		return nil, err
	}
	return promotions, nil
}

func (r *repository) Update(promotion *models.Promotion) error {
	if err := r.db.Save(promotion).Error; err != nil {
		return err
	}
	return nil
}

func (r *repository) Delete(id uint) error {
	if err := r.db.Where("id =?", id).Delete(&models.Promotion{}).Error; err != nil {
		return err
	}

	return nil
}

func (r *repository) FindPromotionByProductID(productID uint) (*models.Promotion, error) {
	promotion := &models.Promotion{}
	err := r.db.Where("product_id =?", productID).First(promotion).Error
	if err != nil {
		return nil, err
	}
	return promotion, nil
}

func (r *repository) FindPromotionByCode(code string) (*models.Promotion, error) {
	promotion := &models.Promotion{}
	err := r.db.Preload("Product").Where("code =?", code).First(promotion).Error
	if err != nil {
		return nil, err
	}
	return promotion, nil
}

func (r *repository) FindPromotionByID(promotionID uint) (*models.Promotion, error) {
	promotion := &models.Promotion{}
	err := r.db.Where("id = ?", promotionID).First(promotion).Error
	if err != nil {
		return nil, err
	}
	return promotion, nil
}

func (r *repository) DeletePromotionID(promotionID uint) error {
	if err := r.db.Model(&models.Cart{}).
		Where("promotion_id =?", promotionID).
		Update("promotion_id", nil).Error; err != nil {
		return err
	}

	return nil
}
