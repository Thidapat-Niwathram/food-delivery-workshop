package main

import (
	"food-delivery-workshop/internal/core/database"
	cart "food-delivery-workshop/internal/pkg/cart"
	"food-delivery-workshop/internal/pkg/product"
	"food-delivery-workshop/internal/pkg/promotion"
	"food-delivery-workshop/internal/pkg/user"
	"log"
	routes "food-delivery-workshop/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	database.ConnectDB()

	userRepository := user.NewRepository(database.DB)
	userService := user.NewService(userRepository)
	productRepository := product.NewRepository(database.DB)
	productService := product.NewService(productRepository)
	promotionRepository := promotion.NewRepository(database.DB)
	promotionService := promotion.NewService(promotionRepository)
	cartRepository := cart.NewRepository(database.DB)
	cartService := cart.NewService(cartRepository,promotionRepository, productRepository)

	app := fiber.New()

	routes.SetupRoutes(app, userService, productService, cartService, promotionService)


	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}

}
