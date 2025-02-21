package promotion

import (
	"errors"
	"food-delivery-workshop/internal/models"
	product "food-delivery-workshop/internal/pkg/product"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	cart "food-delivery-workshop/internal/pkg/cart"
)

type Service interface {
	Create(c *fiber.Ctx, request *CreateRequest) (*models.Promotion, error)
	Update(c *fiber.Ctx, id uint, request *UpdateRequest) (*models.Promotion, error)
	GetByID(id uint) (*models.Promotion, error)
	GetAll() ([]*models.Promotion, error)
	Delete(c *fiber.Ctx, id uint) error
}

type service struct {
	repo        Repository
	productRepo product.Repository
	cartRepo    cart.Repository
}

func NewService(repo Repository, productRepo product.Repository, cartRepo cart.Repository) Service {
	return &service{repo: repo, productRepo: productRepo, cartRepo: cartRepo}
}

func (s *service) Create(c *fiber.Ctx, request *CreateRequest) (*models.Promotion, error) {
	promoCode, err := s.cartRepo.FindPromotionByCode(request.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("find promotion error: %v", err)
		return nil, err
	}

	if promoCode != nil {
        logrus.Errorf("promotion is already exist: %v", err)
        return nil, errors.New("promotion is already exist")
    }

	existingPromotion, err := s.repo.FindPromotionByProductID(request.ProductID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("find promotion error: %v", err)
		return nil, err
	}

	if existingPromotion != nil {
		logrus.Errorf("promotion is already exist: %v", err)
        return nil, errors.New("promotion is already exist")
    }

	promotion := &models.Promotion{}
	_ = copier.Copy(promotion, request)
	if err := s.repo.Create(promotion); err != nil {
		logrus.Errorf("create promotion error: %v", err)
		return nil, err
	}

	s.repo.Preload(promotion)
	return promotion, nil
}

func (s *service) Update(c *fiber.Ctx, id uint, request *UpdateRequest) (*models.Promotion, error) {
	promotion := &models.Promotion{}
	if err := s.repo.FindByID(id, promotion); err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("find promotion error: %v", err)
        return nil, err		
	}

	promoCode, err := s.cartRepo.FindPromotionByCode(request.Code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("find promotion error: %v", err)
		return nil, err
	}

	if promoCode != nil && promoCode.ID!= id {
        logrus.Errorf("promotion is already exist: %v", err)
        return nil, errors.New("promotion is already exist")
    }

	existingPromotion, err := s.repo.FindPromotionByProductID(request.ProductID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
        logrus.Errorf("find promotion error: %v", err)
        return nil, err
    }
	if existingPromotion != nil && existingPromotion.ID != id {
		logrus.Errorf("promotion is already exist: %v", request.ProductID)
        return nil, errors.New("promotion is already exist")
    }

	_ = copier.Copy(promotion, request)
    if err := s.repo.Update(promotion); err != nil {
        logrus.Errorf("update promotion error: %v", err)
        return nil, err
    }

	s.repo.Preload(promotion)
	return promotion, nil
}

func (s *service) GetByID(id uint) (*models.Promotion, error) {
	promotion := &models.Promotion{}
	if err := s.repo.FindByID(id, promotion); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("promotion not found")
		}
		logrus.Errorf("find promotion error: %v", err)
		return nil, err
	}

	return promotion, nil
}

func (s *service) GetAll() ([]*models.Promotion, error) {
	promotions, err := s.repo.FindAll()
	if err != nil {
		logrus.Errorf("find all promotion error: %v", err)
		return nil, err
	}

	return promotions, nil
}

func (s *service) Delete(c *fiber.Ctx, id uint) error {
	promotion := &models.Promotion{}
	if err := s.repo.FindByID(id, promotion); err != nil {
		logrus.Errorf("find promotion error: %v", err)
		return err
	}

	err := s.repo.Delete(id)
	if err != nil {
		logrus.Errorf("delete promotion error: %v", err)
		return err
	}
	return nil
}
