package user

import (
	"errors"
	"food-delivery-workshop/internal/get"
	"food-delivery-workshop/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(c *fiber.Ctx, request *CreateRequest) error
	Login(request *LoginRequest) (*models.User, error)
	GetUserByID(request get.GetOne[uint]) (*models.User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(c *fiber.Ctx, request *CreateRequest) error {
	if err := s.validateRegister(request); err != nil {
		logrus.Errorf("validate register error: %v", err)
		return err
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Errorf("hash password error: %v", err)
		return err
	}
	request.Password = string(hashPassword)
	user := &models.User{}
	_ = copier.Copy(user, request)
	err = s.repo.Create(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Login(request *LoginRequest) (*models.User, error) {
	user := &models.User{}
	if err := s.repo.FindByEmail(request.Email, user); err != nil {
		logrus.Errorf("find user by email error: %v", err)
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		logrus.Warnf("compare password error: %v", err)
		return nil, errors.New("invalid email or password")
	}
	user.Password = ""
	return user, nil
}

func (s * service) GetUserByID(request get.GetOne[uint]) (*models.User, error) {
	user := &models.User{}
	if err := s.repo.FindByID(request.GetID(), user); err != nil{
		logrus.Errorf("find user by id error: %v", err)
		return nil, errors.New("user not found")
	}

	return user, nil
}