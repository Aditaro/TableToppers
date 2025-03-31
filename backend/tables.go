package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/supabase-go"
)

// Table struct
type Table struct {
	ID           string `json:"id"`
	RestaurantID string `json:"restaurant_id"`
	Number       int    `json:"number"`
	Seats        int    `json:"seats"`
	Status       string `json:"status"`
}

// SetupTableRoutes registers the table routes
func SetupTableRoutes(app *fiber.App, client *supabase.Client) {
	app.Get("/restaurants/:restaurantId/tables", func(c *fiber.Ctx) error {
		restaurantID := c.Params("restaurantId")

		var tables []Table
		resp, _, err := client.From("tables").Select("*").Eq("restaurant_id", restaurantID).Execute()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tables"})
		}

		if err := resp.Unmarshal(&tables); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to parse response"})
		}

		return c.JSON(tables)
	})

	app.Post("/restaurants/:restaurantId/tables", func(c *fiber.Ctx) error {
		restaurantID := c.Params("restaurantId")
		var table Table

		if err := c.BodyParser(&table); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		table.RestaurantID = restaurantID

		_, _, err := client.From("tables").Insert(table).Execute()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create table"})
		}

		return c.Status(201).JSON(fiber.Map{"message": "Table created successfully"})
	})
}
