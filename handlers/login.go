package handlers

import (
	"context"
	"prescription/db"
	"prescription/models"
	"prescription/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// CreateUser godoc
// @Summary Create a new user
// @Description Admin can create new users
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.UserDetails true "User data"
// @Success 200 {object} models.UserDetails
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /register [post]
func Register(fi *fiber.Ctx) error {
	var user models.UserDetails

	if err := fi.BodyParser(&user); err != nil {
		return fi.JSON("invalid input")
	}

	user.Role = strings.Title(strings.ToLower(user.Role))

	var adminExists bool
	err := db.Postdb.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM users WHERE role='Admin')`).Scan(&adminExists)
	if err != nil {
		return fi.JSON(err.Error())
	}

	if user.Role == "Admin" && adminExists {
		return fi.JSON("Admin already exists")
	}

	var exists bool
	err = db.Postdb.QueryRow(context.Background(),
		`SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, user.Email).Scan(&exists)
	if err != nil {
		return fi.JSON(err.Error())
	}
	if exists {
		return fi.JSON("User already exists")
	}

	hashpassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fi.JSON("Cannot hash password")
	}
	user.Password = string(hashpassword)

	err = db.Postdb.QueryRow(context.Background(),
		`INSERT INTO users (name,email,password,role) VALUES ($1,$2,$3,$4) RETURNING id,name,email,password,role`,
		user.Name, user.Email, user.Password, user.Role).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.Role)
	if err != nil {
		return fi.JSON(err.Error())
	}

	return fi.JSON(user)
}

// Login godoc
// @Summary Login user
// @Description Authenticate user with email & password, returns JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body models.LoginInput true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /login [post]
func Login(fi *fiber.Ctx) error {
	var input models.LoginInput
	err := fi.BodyParser(&input)
	if err != nil {
		return fi.JSON("invalid input")
	}

	var user models.UserDetails
	err = db.Postdb.QueryRow(context.Background(), `SELECT id, name, email, role, password FROM users WHERE email=$1`, input.Email).Scan(&user.Id, &user.Name, &user.Email, &user.Role, &user.Password)
	if err != nil {
		return fi.JSON("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return fi.JSON("invalid password")
	}

	token, err := utils.GenerateToken(user.Id, user.Role)
	if err != nil {
		return fi.JSON("token not generated")
	}

	return fi.JSON(fiber.Map{"token": token, "user": user})

}
