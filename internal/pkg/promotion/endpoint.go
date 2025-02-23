package promotion

import (
	"food-delivery-workshop/internal/get"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

// Create Create promotion
// @Summary Create a promotion
// @Description Create a promotion
// @Tags promotion
// @Accept  json
// @Produce  json
// @Param request body CreateRequest true "request body"
// @Success 200 {object} models.Promotion
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /promotions [post]
func Create(c *fiber.Ctx, service Service) error {
	request := new(CreateRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := validatePromotionReq(request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	promotion, err := service.Create(c, request)
	if err != nil {
		if err.Error() == "promotion is already exist" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(promotion)
}

// Update update promotion
// @Summary update a promotion
// @Description update a promotion
// @Tags promotion
// @Accept  json
// @Produce  json
// @Param id path uint true "Promotion ID"
// @Param request body Request true "request data"
// @Success 200 {object} models.Promotion
// @Failure 401 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /promotions/{id} [put]
func Update(c *fiber.Ctx, service Service) error {
	promotionIDstr := c.Params("id")
	promotionID, err := strconv.ParseUint(promotionIDstr, 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	request := new(UpdateRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validatePromotionReq(request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}


	request.ID = uint(promotionID)
	updatedPromotion, err := service.Update(c, request)
	if err != nil {
		if err.Error() == "promotion not found" || err.Error() == "promotion is already exist" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(updatedPromotion)
}

// @Summary Delete a promotion
// @Description Delete a promotion by ID
// @Tags promotion
// @Accept json
// @Produce json
// @Param id path int true "Promotion ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /promotions/{id} [delete]
func Delete(c *fiber.Ctx, service Service) error {
	promotionIDstr := c.Params("id")
	promotionID, err := strconv.ParseUint(promotionIDstr, 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid promotion ID",
		})
	}

	promotionIDUint := uint(promotionID)
	ID := &get.GetOne[uint]{ID: promotionIDUint}
	err = service.Delete(c, ID)
	if err != nil {
		if err.Error() == "Promotion not found" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Promotion deleted successfully"})
}

// @Summary Get all Promotions
// @Description Get all Promotions
// @Tags promotion
// @Accept json
// @Produce json
// @Success 200 {array} models.Promotion
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /promotions [get]
func GetAllPromotion(c *fiber.Ctx, service Service) error {
	promotions, err := service.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error getting Promotions",
		})
	}
	return c.Status(fiber.StatusOK).JSON(promotions)
}

// @Summary Get promotion by id
// @Description Get promotion by id
// @Tags promotion
// @Accept json
// @Produce json
// @Param id path int true "promotion ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /promotions/{id} [get]
func GetPromotionByID(c *fiber.Ctx, service Service) error {
	promotionIDstr := c.Params("id")
	promotionID, err := strconv.ParseUint(promotionIDstr, 10, 32)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid promotion ID",
		})
	}

	request := &get.GetOne[uint]{
		ID: uint(promotionID),
	}
	promotion, err := service.GetByID(request)
	if err != nil {
		if err.Error() == "Promotion not found" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{
				"error": "Promotion not found",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error getting promotion",
		})
	}
	return c.Status(http.StatusOK).JSON(promotion)
}
