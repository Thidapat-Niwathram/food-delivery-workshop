package promotion

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func validatePromotionReq(request interface{}) error {
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
		logrus.Errorf("error validate promotion request: %v", err)
		return err
	}
	return nil
}