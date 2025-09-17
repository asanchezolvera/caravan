package products

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"caravan/internal/models"
)

// GetProducts handles the GET /products request and returns a list of products.
func GetProducts(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query(context.Background(), "SELECT * FROM products")
		if err != nil {
			log.Printf("Error querying products: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
		defer rows.Close()

		var products []models.Product
		for rows.Next() {
			var p models.Product
			if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.CreatedAt, &p.UpdatedAt); err != nil {
				log.Printf("Error scanning product row: %v", err)
				continue
			}
			products = append(products, p)
		}

		if rows.Err() != nil {
			log.Printf("Error iterating over product rows: %v", rows.Err())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process products",
			})
		}

		return c.Status(fiber.StatusOK).JSON(products)
	}
}
