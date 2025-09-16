package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/argon2"
)

// User represents the data for a new user registration.
type User struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var db *pgxpool.Pool

func initDB() {
	var err error
	connStr := "postgresql://asanchezolvera:nKE0wa8Lowb676hKTF@localhost:5432/dev_db"
	db, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database not reachable: %v", err)
	}
	log.Println("Successfully connected to database!")
}

func main() {
	initDB()
	defer db.Close()

	app := fiber.New()

	// Endpoint for user registration.
	app.Post("/register", func(c *fiber.Ctx) error {
		// Parse the request body into a User struct.
		user := new(User)
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Cannot parse JSON",
			})
		}

		// Hash the password using argon2.
		hash := argon2.IDKey([]byte(user.Password), []byte("A_RANDOM_SALT"), 1, 64*1024, 4, 32)
		hashedPassword := fmt.Sprintf("%x", hash)

		// Insert the user into the database.
		_, err := db.Exec(context.Background(), "INSERT INTO users (email, password_hash) VALUES ($1, $2)", user.Email, string(hashedPassword))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Could not create user",
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "User registered successfully",
		})
	})

	// Endpoint for user login.
	app.Post("/login", func(c *fiber.Ctx) error {
		return c.SendString("User login endpoint")
	})

	log.Fatal(app.Listen(":3000"))
}
