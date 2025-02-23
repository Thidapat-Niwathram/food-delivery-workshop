package cart

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
)

// Create Add item to cart
// @Summary Create a Cart
// @Description Create a Cart
// @Tags cart
// @Accept  json
// @Produce  json
// @Param request body CreateRequest true "Cart  request"
// @Success 201 {object} models.Cart
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /cart [post]
func Create(c *fiber.Ctx, service Service) error {
	userID := c.Locals("user_id").(float64)
	request := new(CreateRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	request.UserID = uint(userID)
	cartItem, err := service.Create(c, request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(cartItem)
}

// Update updates cart
// @Summary Update a Cart
// @Description Update a Cart
// @Tags cart
// @Accept  json
// @Produce  json
// @Param request body UpdateRequest true "Cart Request"
// @Success 200 {object} models.Cart
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /cart [Put]
func Update(c *fiber.Ctx, service Service) error {
	userID := c.Locals("user_id").(float64)
	request := new(UpdateRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	request.UserID = uint(userID)
	updateCart, err := service.Update(c, request)
	if err != nil {
		if err.Error() == "cart not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(updateCart)
}

// ApplyPromotion apply promotion
// @Summary Apply promotion 
// @Description Apply a promotion code 
// @Tags cart
// @Accept json
// @Produce json
// @Param request body PromotionRequest true "Promotion code request"
// @Success 200 {object} models.Cart
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /cart/promotion [post]
func ApplyPromotion(c *fiber.Ctx, service Service) error {
	userID := c.Locals("user_id").(float64)
	request := new(PromotionRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	request.UserID = uint(userID)
	cart, err := service.ApplyPromotion(c, request)
	if err != nil {
		if err.Error() == "cart not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(cart)
}

// RemoveCartItem remove cart item
// @Summary Remove an item from the cart
// @Description Delete 
// @Tags cart
// @Accept json
// @Produce json
// @Param product_id path int true "Product ID to remove"
// @Success 200 {object} models.Cart "Updated cart details"
// @Failure 400 {object} map[string]string 
// @Failure 401 {object} map[string]string 
// @Failure 500 {object} map[string]string 
// @Router /cart/item/{product_id} [delete]
// @Security ApiKeyAuth
func RemoveCartItem(c *fiber.Ctx, service Service) error {
	userID := c.Locals("user_id").(float64)
	productID, err := strconv.Atoi(c.Params("product_id"))
	if err != nil || productID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid product_id",
		})
	}

	updatedCart, err := service.RemoveItem(c, &RemoveItemRequest{UserID: uint(userID), ProductID: uint(productID)})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedCart)
}

// GetCart get cart
// @Summary Get all cart items 
// @Description Get all cart items 
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} models.Cart "Cart"
// @Failure 400 {object} map[string]string 
// @Failure 401 {object} map[string]string 
// @Failure 500 {object} map[string]string 
// @Router /cart [get]
// @Security ApiKeyAuth
func GetAllCart(c *fiber.Ctx, service Service) error {
    userID := c.Locals("user_id").(float64)
    cart, err := service.GetAllCart(c, &GetAllRequests{UserID: uint(userID)})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return c.Status(fiber.StatusOK).JSON(cart)
}
