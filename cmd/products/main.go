package main

import (
	"context"
	"fmt"
	"log"
	"os"

	handlers "caravan/internal/products"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	_ "caravan/internal/models"
)

var db *pgxpool.Pool

func initDB() {
	var err error
	connStr := fmt.Sprintf(
		"postgres://%s:%s:@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

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
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	initDB()
	defer db.Close()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello from the products microservice!")
	})

	app.Get("/products", handlers.GetProducts(db))

	log.Fatal(app.Listen(":3001"))
}
