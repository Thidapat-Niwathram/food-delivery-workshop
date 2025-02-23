package cart

import (
	"errors"
	"food-delivery-workshop/internal/models"
	product "food-delivery-workshop/internal/pkg/product"
	promotion "food-delivery-workshop/internal/pkg/promotion"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Service interface {
	Create(c *fiber.Ctx, request *CreateRequest) (*models.Cart, error)
	ApplyPromotion(c *fiber.Ctx, request *PromotionRequest) (*models.Cart, error)
	Update(c *fiber.Ctx, request *UpdateRequest) (*models.Cart, error)
	RemoveItem(c *fiber.Ctx, request *RemoveItemRequest) (*models.Cart, error)
	GetAllCart(c *fiber.Ctx, request *GetAllRequests) (*models.Cart, error)
}

type service struct {
	repo        Repository
	promoRepo   promotion.Repository
	productRepo product.Repository
}

func NewService(repo Repository, promoRepo promotion.Repository, productRepo product.Repository) Service {
	return &service{repo: repo, promoRepo: promoRepo, productRepo: productRepo}
}

func (s *service) CalculateCartItem(cartItem *models.CartItem) error {
	product, err := s.productRepo.FindByProductID(cartItem.ProductID)
	if err != nil {
		logrus.Errorf("find product error: %v", err)
		return err
	}

	cartItem.Price = product.Price
	cartItem.TotalPrice = cartItem.Price * float64(cartItem.Quantity)

	return nil
}

func (s *service) CalculateCart(cart *models.Cart) error {
	var totalAmount float64
	for _, cartItem := range cart.CartItems {
		if err := s.CalculateCartItem(cartItem); err != nil {
			logrus.Errorf("calculate cart item error: %v", err)
			return err
		}
		totalAmount += cartItem.TotalPrice
	}

	cart.SubTotal = totalAmount
	cart.Total = totalAmount - cart.Discount

	return nil
}
func (s *service) Create(c *fiber.Ctx, request *CreateRequest) (*models.Cart, error) {
	existingCart, err := s.repo.FindCartByUserID(request.UserID)
	if existingCart != nil {
		return nil, errors.New("cart already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("find cart error: %v", err)
		return nil, err
	}

	if len(request.CartItemRequests) == 0 {
		return nil, errors.New("cart_items cannot be empty")
	}

	cart := &models.Cart{
		UserID: request.UserID,
	}
	if err := s.repo.CreateCart(cart); err != nil {
		logrus.Errorf("create cart error: %v", err)
		return nil, err
	}

	cartItems := []*models.CartItem{}
	for _, req := range request.CartItemRequests {
		if req.Quantity <= 0 {
			return nil, errors.New("quantity must be more than 0")
		}

		cartItem := &models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}

		if err := s.CalculateCartItem(cartItem); err != nil {
			logrus.Errorf("calculate cart item error: %v", err)
			return nil, err
		}

		cartItems = append(cartItems, cartItem)
	}

	for _, item := range cartItems {
		if err := s.repo.CreateCartItem(item); err != nil {
			logrus.Errorf("crate cart item error: %v", err)
			return nil, err
		}
	}

	cart.CartItems = cartItems

	if err := s.CalculateCart(cart); err != nil {
		logrus.Errorf("calculate cart error: %v", err)
		return nil, err
	}

	if err := s.repo.UpdateCart(cart); err != nil {
		logrus.Errorf("update cart error: %v", err)
		return nil, err
	}

	s.repo.Preload(cart)
	return cart, nil
}

func (s *service) Update(c *fiber.Ctx, request *UpdateRequest) (*models.Cart, error) {
	cart, err := s.repo.FindCartByUserID(request.UserID)
	if err != nil {
		logrus.Errorf("find cart error: %v", err)
		return nil, errors.New("cart not found")
	}

	if len(request.CartItemRequests) == 0 {
		return nil, errors.New("cart_items cannot be empty")
	}

	if err := s.repo.DeleteAllCartItems(cart.ID); err != nil {
		logrus.Errorf("delete all cart items error: %v", err)
		return nil, err
	}

	cartItems := []*models.CartItem{}
	for _, req := range request.CartItemRequests {
		if req.Quantity <= 0 {
			// ถ้าจำนวนเป็น 0 หรือ น้อยกว่า ให้ข้ามไป (ไม่เพิ่มสินค้านี้)
			continue
		}

		cartItem := &models.CartItem{
			CartID:    cart.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}

		if err := s.CalculateCartItem(cartItem); err != nil {
			logrus.Errorf("calculate cart item error: %v", err)
			return nil, err
		}

		cartItems = append(cartItems, cartItem)
	}

	for _, item := range cartItems {
		if err := s.repo.CreateCartItem(item); err != nil {
			logrus.Errorf("crate cart item error: %v", err)
			return nil, err
		}
	}

	cart.CartItems = cartItems

	if err := s.CalculateCart(cart); err != nil {
		logrus.Errorf("calculate cart error: %v", err)
		return nil, err
	}

	if err := s.repo.UpdateCart(cart); err != nil {
		logrus.Errorf("update cart error: %v", err)
		return nil, err
	}

	s.repo.Preload(cart)
	return cart, nil
}

func (s *service) ApplyPromotion(c *fiber.Ctx, request *PromotionRequest) (*models.Cart, error) {
	cart, err := s.repo.FindCartByUserID(request.UserID)
	if err != nil {
		logrus.Errorf("find cart error: %v", err)
		return nil, err
	}

	promotion, err := s.promoRepo.FindPromotionByCode(request.PromotionCode)
	if err != nil {
		logrus.Errorf("find promotion error: %v", err)
		return nil, err
	}

	var isValidPromotion bool
	for _, item := range cart.CartItems {
		if err := s.CalculateCartItem(item); err != nil {
			logrus.Errorf("calculate cart item error: %v", err)
			return nil, err
		}

		if item.ProductID == promotion.ProductID {
			isValidPromotion = true
			cart.Discount = promotion.Discount
		}

	}

	if !isValidPromotion {
		return nil, errors.New("promotion is not applicable for items in the cart")
	}

	if err := s.CalculateCart(cart); err != nil {
		logrus.Errorf("calculate cart error: %v", err)
		return nil, err
	}

	cart.PromotionID = &promotion.ID
	if err := s.repo.UpdateCart(cart); err != nil {
		logrus.Errorf("update cart error: %v", err)
		return nil, err
	}

	s.repo.Preload(cart)
	return cart, nil
}

func (s *service) GetAllCart(c *fiber.Ctx, request *GetAllRequests) (*models.Cart, error) {
	cart, err := s.repo.FindCartByUserID(request.UserID)
	if err != nil {
		logrus.Errorf("find cart error: %v", err)
		return nil, err
	}

	var totalAmount float64
	for _, item := range cart.CartItems {
		if item.Product != nil {
			totalAmount += item.TotalPrice
		}
	}

	if cart.PromotionID != nil {
		cart.Discount = cart.Promotion.Discount
	}

	cart.SubTotal = totalAmount
	cart.Total = totalAmount - cart.Discount

	return cart, nil
}

func (s *service) RemoveItem(c *fiber.Ctx, request *RemoveItemRequest) (*models.Cart, error) {
	cart, err := s.repo.FindCartByUserID(request.UserID)
	if err != nil {
		logrus.Errorf("find cart error: %v", err)
		return nil, err
	}
	if _, err := s.repo.FindCartItem(cart.ID, request.ProductID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("cart item not found")
		}
		logrus.Errorf("find cart item error: %v", err)
		return nil, err
	}

	if err := s.repo.DeleteCartItem(cart.ID, request.ProductID); err != nil {
		logrus.Errorf("delete cart item error: %v", err)
		return nil, err
	}

	cartItems, err := s.repo.FindCartItemsByCartID(cart.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("find cart items error: %v", err)
		return nil, err
	}

	var newSubTotal float64
	for _, item := range cartItems {
		newSubTotal += item.TotalPrice
	}
	cart.SubTotal = newSubTotal
	cart.Total = newSubTotal

	remainingItems, err := s.repo.CountCartItems(cart.ID)
	if err != nil {
		logrus.Errorf("count cart items error: %v", err)
		return nil, err
	}

	if remainingItems == 0 {
		if err := s.repo.DeleteCart(cart.ID); err != nil {
			logrus.Errorf("delete cart error: %v", err)
			return nil, err
		}
		return nil, nil
	}

	if err := s.repo.Preload(cart); err != nil {
		logrus.Errorf("preload cart error: %v", err)
		return nil, err
	}
	
	return cart, nil
}
