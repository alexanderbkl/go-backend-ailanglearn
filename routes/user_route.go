package routes

import (
	"go-backend-ailanglearn/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRoute(app *fiber.App) {


		app.Get("/user/:userId", controllers.GetAUser)
		app.Put("/user/:userId", controllers.EditAUser)
		app.Delete("/user/:userId", controllers.DeleteAUser)
		app.Get("/users", controllers.GetAllUsers)

		app.Post("/signup", controllers.HandleSignup)
		app.Post("/signin", controllers.HandleSignin)
		app.Get("/getprofile", controllers.GetProfile)
		app.Post("/createnote", controllers.CreateNote)
	
}
