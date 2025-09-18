package main

import (
	"fmt"
	"log"
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"caravan/internal/auth"
	"caravan/internal/models"
)

var db *gorm.DB

func initDB() {
	var err error

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("Unable to load environment variables: %v", err)
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = fmt.Sprintf(
			"postgres://%s:%s:@%s:%s/%s",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)
	}

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Successfully connected to database!")

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Successfully migrated database!")
}

func main() {

	initDB()

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

	app.Get("/profile", auth.GetUserProfile(service))

	log.Fatal(app.Listen(":3000"))
}
