package product

import (
	"errors"
	"food-delivery-workshop/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Service interface {
	Create(c *fiber.Ctx, request *CreateRequest) (*models.Product, error)
	Update(c *fiber.Ctx, productID uint, request *UpdateRequest) (*models.Product, error)
	GetProductByID(id uint) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
	Delete(c *fiber.Ctx, id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

// Create create a product
func (s *service) Create(c *fiber.Ctx, request *CreateRequest) (*models.Product, error) {
	product := &models.Product{}
	productName, err := s.repo.FindByProductName(request.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("find product name error: %v", err)
		return nil, err
	}

	if productName != nil {
		logrus.Errorf("product name already exist: %v", err)
		return nil, errors.New("product name already exist")
	}

	_ = copier.Copy(product, request)
	if err := s.repo.Create(product); err != nil {
		logrus.Errorf("create product error: %v", err)
		return nil, err
	}

	return product, nil
}

func (s *service) Update(c *fiber.Ctx, id uint, request *UpdateRequest) (*models.Product, error) {
	product := &models.Product{}
	if err := s.repo.FindByID(id, product); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		logrus.Errorf("find product by id error: %v", err)
		return nil, err
	}

	productName, err := s.repo.FindByProductName(request.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logrus.Errorf("find product name error: %v", err)
		return nil, err
	}
	if productName != nil {
		logrus.Errorf("product name already exist: %v", err)
		return nil, errors.New("product name already exist")
	}
	_ = copier.Copy(product, request)
	if err := s.repo.Update(product); err != nil {
		logrus.Errorf("update product error: %v", err)
		return nil, err
	}

	return product, nil
}

func (s *service) GetProductByID(id uint) (*models.Product, error) {
	product := &models.Product{}
	if err := s.repo.FindByID(id, product); err != nil {
		return nil, err
	}
	//s.repo.Preload(product)
	return product, nil
}

func (s *service) GetAllProducts() ([]models.Product, error) {
	products, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *service) Delete(c *fiber.Ctx, id uint) error {
	product := &models.Product{}
	if err := s.repo.FindByID(id, product); err != nil {
		return errors.New("product not found")
	}

	err := s.repo.Delete(id)
	if err != nil {
		logrus.Errorf("delete product error: %v", err)
		return err
	}

	return nil
}
