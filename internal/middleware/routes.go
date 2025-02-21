package middleware

import (
	"food-delivery-workshop/internal/pkg/cart"
	"food-delivery-workshop/internal/pkg/product"
	"food-delivery-workshop/internal/pkg/promotion"
	"food-delivery-workshop/internal/pkg/user"
	"food-delivery-workshop/internal/auth"

	fiberSwagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	jwtware "github.com/gofiber/jwt/v3"
	jwt "github.com/golang-jwt/jwt/v4"
)

func SetupRoutes(app *fiber.App, userService user.Service, productService product.Service, cartService cart.Service, promotionService promotion.Service) {
	users := make(map[string]string)
	users["username"] = "password"
	basicAuth := basicauth.New(basicauth.Config{
		Users: users,
		Unauthorized: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		},
	})

	auth := jwtware.New(jwtware.Config{
		SigningKey: []byte(auth.SecretKey),
		ErrorHandler: func(c *fiber.Ctx, _ error) error {
			return basicAuth(c)
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)
			c.Locals("user_id", claims["user_id"])
			return c.Next()
		},
	})

	// Routes for Users
	app.Post("/users/login", func(c *fiber.Ctx) error {
		return user.Login(c, userService)
	})
	app.Post("/users/register", func(c *fiber.Ctx) error {
		return user.Register(c, userService)
	})
	app.Get("/me", auth, func(c *fiber.Ctx) error {
		return user.GetUserByID(c, userService)
	})

	// Routes for Products
	app.Post("/products", auth, func(c *fiber.Ctx) error {
		return product.Create(c, productService)
	})
	app.Put("/products/:id", auth, func(c *fiber.Ctx) error {
		return product.Update(c, productService)
	})
	app.Delete("/products/:id", auth, func(c *fiber.Ctx) error {
		return product.Delete(c, productService)
	})
	app.Get("/products", auth, func(c *fiber.Ctx) error {
		return product.GetAllProduct(c, productService)
	})
	app.Get("/products/:id", auth, func(c *fiber.Ctx) error {
		return product.GetProductByID(c, productService)
	})

	// Routes for Promotions
	app.Post("/promotions", auth, func(c *fiber.Ctx) error {
		return promotion.Create(c, promotionService)
	})
	app.Put("/promotions/:id", auth, func(c *fiber.Ctx) error {
		return promotion.Update(c, promotionService)
	})
	app.Delete("/promotions/:id", auth, func(c *fiber.Ctx) error {
		return promotion.Delete(c, promotionService)
	})
	app.Get("/promotions", auth, func(c *fiber.Ctx) error {
		return promotion.GetAllPromotion(c, promotionService)
	})
	app.Get("/promotions/:id", auth, func(c *fiber.Ctx) error {
		return promotion.GetPromotionByID(c, promotionService)
	})

	// Routes for Cart
	app.Post("/cart", auth, func(c *fiber.Ctx) error {
		return cart.Create(c, cartService)
	})
	app.Put("/cart", auth, func(c *fiber.Ctx) error {
		return cart.Update(c, cartService)
	})
	app.Post("/cart/promotion", auth, func(c *fiber.Ctx) error {
		return cart.ApplyPromotion(c, cartService)
	})
	app.Delete("/cart/item/:product_id", auth, func(c *fiber.Ctx) error {
		return cart.RemoveCartItem(c, cartService)
	})
	app.Get("/cart", auth, func(c *fiber.Ctx) error {
		return cart.GetAllCart(c, cartService)
	})

	// Swagger Route
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)
}
