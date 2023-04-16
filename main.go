package main

import (
	"go-backend-ailanglearn/configs"
	"go-backend-ailanglearn/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {

	app := fiber.New()
	app.Use(cors.New())

	configs.ConnectDB()


	routes.UserRoute(app)

	//expose host for other devices
	app.Listen(":3001")

	
}