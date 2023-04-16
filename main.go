package main

import (
	"go-backend-ailanglearn/configs"
	"go-backend-ailanglearn/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {

	app := fiber.New()
	app.Use(cors.New())

	app.Use(limiter.New(limiter.Config{
		Max:               20,
		KeyGenerator: 	func(c *fiber.Ctx) string { return c.IP() },
		Expiration:        30 * time.Second,
		LimiterMiddleware: limiter.SlidingWindow{},
	}))

	configs.ConnectDB()

	routes.UserRoute(app)

	//expose host for other devices
	app.Listen(":3001")

}
