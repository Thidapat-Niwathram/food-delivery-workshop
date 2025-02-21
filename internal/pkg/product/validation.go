package product

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func validateProductReq(request interface{}) error {
	validate := validator.New()
	err := validate.Struct(request)
	if err != nil {
        logrus.Error("error validate product request: %v", err)
        return err
    }
	return nil
}