package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/supabase-go"
)

// Reservation struct
type Reservation struct {
	ID           string `json:"id"`
	RestaurantID string `json:"restaurant_id"`
	UserID       string `json:"user_id"`
	Date         string `json:"date"`
	Time         string `json:"time"`
	Guests       int    `json:"guests"`
	Status       string `json:"status"`
}

// SetupReservationsRoutes registers the reservation routes
func SetupReservationsRoutes(app *fiber.App, client *supabase.Client) {
	app.Get("/restaurants/:restaurantId/reservations", func(c *fiber.Ctx) error {
		restaurantID := c.Params("restaurantId")
		date := c.Query("date")

		query := client.From("reservations").Select("*").Eq("restaurant_id", restaurantID)
		if date != "" {
			query = query.Eq("date", date)
		}

		var reservations []Reservation
		resp, _, err := query.Execute()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch reservations"})
		}

		if err := resp.Unmarshal(&reservations); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to parse response"})
		}

		return c.JSON(reservations)
	})

	app.Post("/restaurants/:restaurantId/reservations", func(c *fiber.Ctx) error {
		restaurantID := c.Params("restaurantId")
		var reservation Reservation

		if err := c.BodyParser(&reservation); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
		}

		reservation.RestaurantID = restaurantID

		_, _, err := client.From("reservations").Insert(reservation).Execute()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create reservation"})
		}

		return c.Status(201).JSON(fiber.Map{"message": "Reservation created successfully"})
	})
