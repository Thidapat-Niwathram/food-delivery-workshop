package cart

import (
	"food-delivery-workshop/internal/models"
	"gorm.io/gorm"
	"github.com/sirupsen/logrus"

)

type Repository interface {
	CreateCartItem(cartItem *models.CartItem) error
	FindCartByUserID(userID uint) (*models.Cart, error)
	FindByProductID(productID uint) (*models.Product, error)
	CreateCart(cart *models.Cart) error
	UpdateCart(cart *models.Cart) error
	UpdateCartItem(cart *models.CartItem) error
	FindPromotionByCode(code string) (*models.Promotion, error)
	FindPromotionByID(promotionID uint) (*models.Promotion, error)
	FindCartItem(cartID uint, productID uint) (*models.CartItem, error)
	DeleteCartItem(cartID uint, productID uint) error
	RemoveItem(cartID uint, cartItemID uint) error
	DeleteCart(cartID uint) error 
	Preload(cart interface{}) error
	DeleteAllCartItems(cartID uint) error
	CountCartItems(cartID uint) (int64, error)
	FindCartItemsByCartID(cartID uint) ([]*models.CartItem, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Preload(cart interface{}) error {
	return r.db.Preload("CartItems.Product").Find(cart).Error
}
func (r *repository) CreateCartItem(cartItem *models.CartItem) error {
	if err := r.db.Create(cartItem).Error; err != nil {
		return err
	}
	return r.db.Preload("Product").First(cartItem, cartItem.ID).Error
}

func (r *repository) FindCartByUserID(userID uint) (*models.Cart, error) {
	cart := &models.Cart{}
	err := r.db.Preload("CartItems.Product").Where("user_id = ?", userID).First(cart).Error
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (r *repository) FindByProductID(productID uint) (*models.Product, error) {
	product := &models.Product{}
	err := r.db.
		Preload("Promotion").
		Where("id =?", productID).
		First(product).Error
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *repository) FindPromotionByCode(code string) (*models.Promotion, error) {
	promotion := &models.Promotion{}
    err := r.db.Preload("Product").Where("code =?", code).First(promotion).Error
    if err != nil {
        return nil, err
    }
    return promotion, nil
}

func (r *repository) CreateCart(cart *models.Cart) error {
	err := r.db.Create(cart).Error
	if err != nil {
		return err
	}

	return r.db.Preload("CartItems.Product").First(cart, cart.ID).Error
}

func (r *repository) UpdateCart(cart *models.Cart) error {
	err := r.db.Save(cart).Error
	if err != nil {
		logrus.Errorf("failed to update cart: %v", err)
		return err
	}
	return nil
}

func (r *repository) UpdateCartItem(cart *models.CartItem) error {
	err := r.db.Save(cart).Error
    if err!= nil {
        logrus.Errorf("failed to update cart item: %v", err)
        return err
    }
    return nil
}

func (r *repository) FindCartItem(cartID uint, productID uint) (*models.CartItem, error) {
	cartItem := &models.CartItem{}
	err := r.db.Where("cart_id = ? AND product_id =?", cartID, productID).First(cartItem).Error
	if err != nil {
		return nil, err
	}

	return cartItem, nil
}

func (r *repository) DeleteCartItem(cartID uint, productID uint) error {
	if err := r.db.Where("cart_id = ? AND product_id = ?", cartID, productID).Delete(&models.CartItem{}).Error; err != nil {
		return err
	}
    return nil
}

func (r *repository) RemoveItem(cartID uint, cartItemID uint) error {
	if err := r.db.Where("card_id = ? AND id = ?", cartID, cartItemID).Delete(&models.CartItem{}).Error; err != nil {
		logrus.Errorf("failed to delete cart item: %v", err)
		return err
	}
	return nil
}

func (r *repository) FindPromotionByID(promotionID uint) (*models.Promotion, error) {
	promotion := &models.Promotion{}
	err := r.db.Where("id = ?", promotionID).First(promotion).Error
	if err != nil {
		return nil, err
	}
	return promotion, nil
}

func (r *repository) DeleteCart(cartID uint) error {
	if err := r.db.Where("id =?", cartID).Delete(&models.Cart{}).Error; err != nil {
        return err
    }
    return nil
}

func (r *repository) DeleteAllCartItems(cartID uint) error {
	if err := r.db.Where("cart_id =?", cartID).Delete(&models.CartItem{}).Error; err != nil {
        return err
    }
    return nil
}

func (r *repository) CountCartItems(cartID uint) (int64, error) {
	var count int64
    err := r.db.Model(&models.CartItem{}).Where("cart_id =?", cartID).Count(&count).Error
    if err != nil {
        return 0, err
    }
    return count, nil
}

func (r *repository) FindCartItemsByCartID(cartID uint) ([]*models.CartItem, error) {
	var cartItems []*models.CartItem
	err := r.db.Where("cart_id =?", cartID).Preload("Product").Find(&cartItems).Error
	if err != nil {
		return nil, err
	}

	return cartItems, nil
}