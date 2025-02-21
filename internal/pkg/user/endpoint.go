package user

import (
	"food-delivery-workshop/internal/auth"
	"net/http"
	"github.com/gofiber/fiber/v2"
)

// RegisterUser register
// @Summary Register a new user
// @Description Registers a new user
// @Tags user
// @Accept  json
// @Produce  json
// @Param user body CreateRequest true "User Registration Data"
// @Success 201 {object} CreateRequest
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/register [post]
func Register(c *fiber.Ctx, service Service) error {
	request := &CreateRequest{}
	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	err := service.Create(c, request)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Server error: " + err.Error(),
		})
	}
	request.Password = ""

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create user success",
		"user":    request,
	})
}

// Login Login user
// @Summary Login
// @Description Logs in a user with email and password
// @Tags user
// @Accept  json
// @Produce  json
// @Param login body LoginRequest true "Login Data"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /users/login [post]
func Login(c *fiber.Ctx, service Service) error {
	request := &LoginRequest{}
	if err := c.BodyParser(request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	user, err := service.Login(request)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	token, err := auth.GenerateJWT(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	return c.JSON(fiber.Map{
		"access_token": token,
	})
}

// GetUserByID Get user by ID
// @Summary Get User Information
// @Description Get User Profile
// @Tags user
// @Accept  json
// @Produce  json
// @Success 200 {object} models.User "User profile"
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security ApiKeyAuth
// @Router /me [get]
func GetUserByID(c *fiber.Ctx, service Service) error {
	authUserID := c.Locals("user_id")
	if authUserID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	authUserIDUint, ok := authUserID.(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid user token",
		})
	}

	userID := uint(authUserIDUint)
	user, err := service.GetUserByID(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	user.Password = ""
	return c.Status(fiber.StatusOK).JSON(user)

}
