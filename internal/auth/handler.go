package auth

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"

	"caravan/internal/models"
)

var JwtSecret = []byte("secret")

// RegisterUser handles the user registration request.
func RegisterUser(svc *Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}
		if err := svc.RegisterUser(c.Context(), user); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to register user",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User registered successfully",
		})
	}
}

// LoginUser handles the user login request.
func LoginUser(svc *Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := new(models.User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		token, err := svc.LoginUser(c.Context(), user)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid email or password",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"token": token,
		})
	}
}

// GetUserProfile handles the protected user profile route.
func GetUserProfile(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(*jwt.MapClaims)
	email := (*claims)["email"].(string)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": fmt.Sprintf("Welcome, %s! You are authorized.", email),
	})
}
