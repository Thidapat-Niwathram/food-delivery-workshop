package user

import (
	"errors"
	"food-delivery-workshop/internal/models"
	"regexp"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (s *service) validateRegister(request *CreateRequest) error {
	if request.FirstName == "" || request.LastName == "" || request.Email == "" || request.Password == "" {
		return errors.New("missing required fields")
	}

	if !isValidEmail(request.Email) {
		return errors.New("invalid email format")
	}

	IsEmailExists, err := s.IsEmailExists(request.Email)
	if err != nil {
		logrus.Errorf("IsEmailExists is error: %v", err)
		return err
	}

	if IsEmailExists {
		return errors.New("email already exists")
	}

	if len(request.Password) < 10 {
		return errors.New("password must be at least 10 characters long")
	}

	if !hasEnglishLetters(request.Password) {
		return errors.New("password must contain at least one english letter")
	}

	if err := validatePhone(request.Phone); err != nil {
		logrus.Errorf("validate phone error: %v", err)
		return err
	}

	if err := validateIDCard(request.IDCard); err != nil {
		logrus.Errorf("validate id card error: %v", err)
		return err
	}

	return nil
}

func (s *service) IsEmailExists(email string) (bool, error) {
	user := &models.User{}
	err := s.repo.FindByEmail(email, user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
            return false, nil 
        }
		logrus.Errorf("error find user by email: %v", err)
		return false, err
	}
	return true, nil
}

func hasEnglishLetters(password string) bool {
	var hasUpper, hasLower bool
	for _, char := range password {
		if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if char >= 'A' && char <= 'Z' {
			hasUpper = true
		}
	}
	return hasUpper && hasLower
}

func isValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func validatePhone(phone string) error {
	phoneRegex := `^0[689]\d{8}$`
    re := regexp.MustCompile(phoneRegex)
    if !re.MatchString(phone) {
        return errors.New("invalid phone number format")
    }
    return nil
}

func validateIDCard(idCard string) error {
	idCardRegex := regexp.MustCompile(`^\d{13}$`)
	if !idCardRegex.MatchString(idCard) {
		return errors.New("invalid ID card format")
	}
	return nil
}
