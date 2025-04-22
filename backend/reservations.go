package main

import (
	"encoding/json" // Import encoding/json

	"github.com/gofiber/fiber/v2"
	"github.com/supabase-community/supabase-go"
	// "log" // Uncomment if you want to add logging back
)

// Reservation struct
type Reservation struct {
	ID           string `json:"id,omitempty"` // omitempty helps during insert
	RestaurantID string `json:"restaurant_id"`
	UserID       string `json:"user_id,omitempty"` // Assuming nullable or set later
	Date         string `json:"date"`
	Time         string `json:"time"`
	Guests       int    `json:"guests"`
	Status       string `json:"status,omitempty"` // Assuming default or set later
}

// SetupReservationsRoutes registers the reservation routes
func SetupReservationsRoutes(app *fiber.App, client *supabase.Client) {
	// Route to get reservations for a specific restaurant, optionally filtered by date
	app.Get("/restaurants/:restaurantId/reservations", func(c *fiber.Ctx) error {
		restaurantID := c.Params("restaurantId")
		date := c.Query("date") // Get optional date query parameter

		// Start building the Supabase query
		// Adjust Select call: Added default count ("") and head (false) arguments
		query := client.From("reservations").Select("*", "", false).Eq("restaurant_id", restaurantID)

		// Add date filter if provided
		if date != "" {
			query = query.Eq("date", date)
		}

		var reservations []Reservation
		// Execute the query
		// Execute now returns []byte, count, error
		respBytes, _, err := query.Execute()
		if err != nil {
			// log.Printf("Error fetching reservations: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch reservations"})
		}

		// Use json.Unmarshal to parse the response bytes
		if err := json.Unmarshal(respBytes, &reservations); err != nil {
			// log.Printf("Error parsing reservation response: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse response"})
		}

		// Return the fetched reservations as JSON
		return c.JSON(reservations)
	})

	// Route to create a new reservation for a specific restaurant
	app.Post("/restaurants/:restaurantId/reservations", func(c *fiber.Ctx) error {
		restaurantID := c.Params("restaurantId")
		var reservation Reservation // Struct to hold the incoming reservation data

		// Parse the request body into the reservation struct
		if err := c.BodyParser(&reservation); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
		}

		// Set the RestaurantID from the URL parameter
		reservation.RestaurantID = restaurantID
		// reservation.Status = "pending" // Example: Set default status

		// Insert the new reservation into the Supabase table
		// Adjust Insert call: Added default upsert (false), returning ("representation"),
		// onConflict (""), and count ("") arguments.
		// Pass the reservation struct (or a slice/map) as the first argument.
		_, _, err := client.From("reservations").Insert(reservation, false, "representation", "", "").Execute()
		if err != nil {
			// log.Printf("Error creating reservation: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create reservation"})
		}

		// Return success message
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Reservation created successfully"})
	})
}

