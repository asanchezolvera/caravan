package main

import (
	"context"
	"fmt"
	"log"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"caravan/internal/auth"
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
	repo := auth.NewRepository(db)
	service := auth.NewService(repo)

	app.Post("/register", auth.RegisterUser(service))
	app.Post("/login", auth.LoginUser(service))

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: auth.JwtSecret,
		},
	}))

	app.Get("/profile", auth.GetUserProfile)

	log.Fatal(app.Listen(":3000"))
}
