package main

import (
	"prescription/db"
	_ "prescription/docs"
	"prescription/handlers"
	"prescription/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"

	"github.com/joho/godotenv"
)

// @title Prescription API
// @version 1.0
// @description Backend service for prescriptions and medicine stock management
// @contact.name Your Name
// @contact.email youremail@example.com
// @host localhost:3000
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	db.ConnectPostDb()
	godotenv.Load()

	app := fiber.New()

	app.Post("/register", middleware.AdminVerification, handlers.Register)
	app.Post("/login", handlers.Login)
	app.Post("/medicine", middleware.AdminVerification, handlers.Medicine)
	app.Delete("/medicine/:medicine_name", middleware.AdminVerification, handlers.RemoveMedicine)
	app.Get("/medicines",middleware.AdminVerification, handlers.GetallMeds)
	app.Post("/updatemeds", middleware.PharmacistVerification, handlers.DispenseStock)
	app.Post("/presc",middleware.DoctorVerification,handlers.MakePresc)

	app.Get("/swagger/*", swagger.HandlerDefault)

	app.Listen(":3000")
}
