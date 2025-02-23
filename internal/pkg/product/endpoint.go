package product

import (
	"food-delivery-workshop/internal/get"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Create Create product
// @Summary Create a product
// @Description Create a product
// @Tags product
// @Accept  json
// @Produce  json
// @Param request body CreateRequest true "Product"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security ApiKeyAuth
// @Router /products [post]
func Create(c *fiber.Ctx, service Service) error {
	request := new(CreateRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := validateProductReq(request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	product, err := service.Create(c, request)
	if err != nil {
		if err.Error() == "product name already exist" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

// Update update product
// @Summary update a product
// @Description update a product
// @Tags product
// @Accept  json
// @Produce  json
// @Param id path uint true "Product ID"
// @Param request body UpdateRequest true "Product data"
// @Success 200 {object} models.Product
// @Failure 401 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Security ApiKeyAuth
// @Router /products/{id} [put]
func Update(c *fiber.Ctx, service Service) error {
	productIDstr := c.Params("id")
	productID, err := strconv.ParseUint(productIDstr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	request := new(UpdateRequest)
	if err := c.BodyParser(request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := validateProductReq(request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	request.ID = uint(productID)
	updatedProduct, err := service.Update(c, request)
	if err != nil {
		if err.Error() == "product not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedProduct)

}

// @Summary Delete a product
// @Description Soft delete a product by ID
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /products/{id} [delete]
func Delete(c *fiber.Ctx, service Service) error {
	productIDstr := c.Params("id")
	productID, err := strconv.ParseUint(productIDstr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}

	productIDUint := uint(productID)
	ID := &get.GetOne[uint]{ID: productIDUint}
	err = service.Delete(c, ID)
	if err != nil {
		if err.Error() == "product not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Product deleted successfully",
	})
}

// @Summary Get all products
// @Description Get all products
// @Tags product
// @Accept json
// @Produce json
// @Success 200 {array} models.Product
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /products [get]
func GetAllProduct(c *fiber.Ctx, service Service) error {
	products, err := service.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error getting products",
		})
	}
	return c.Status(fiber.StatusOK).JSON(products)
}

// @Summary Get product by id
// @Description Get product by id
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security ApiKeyAuth
// @Router /products/{id} [get]
func GetProductByID(c *fiber.Ctx, service Service) error {
	productIDstr := c.Params("id")
	productID, err := strconv.ParseUint(productIDstr, 10, 0)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid product ID",
		})
	}
	productIDUint := uint(productID)
	ID := &get.GetOne[uint]{ID: productIDUint}
	product, err := service.GetProductByID(ID)
	if err != nil {
		if err.Error() == "product not found" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Product not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error getting product",
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}
