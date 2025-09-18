package products

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"caravan/internal/models"
)

// GetProducts handles the GET /products request and returns a list of products.
func GetProducts(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var products []models.Product

		result := db.Find(&products)

		if result.Error != nil {
			log.Printf("Error finding products: %v", result.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process products",
			})
		}

		return c.Status(fiber.StatusOK).JSON(products)
	}
}
