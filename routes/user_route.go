package routes

import (
	"github.com/gofiber/fiber/v2"
	"go-backend-ailanglearn/controllers"
)


func UserRoute(app *fiber.App) {
	app.Get("/user/:userId", controllers.GetAUser)
	app.Put("/user/:userId", controllers.EditAUser)
	app.Delete("/user/:userId", controllers.DeleteAUser)
	app.Get("/users", controllers.GetAllUsers)

	app.Post("/signup", controllers.HandleSignup)
	app.Post("/signin", controllers.HandleSignin)
	app.Post("/createnote", controllers.CreateNote)
}