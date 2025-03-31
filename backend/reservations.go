package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
)

type Reservation struct {
	ID           string `json:"id"`
	RestaurantID string `json:"restaurant_id"`
	UserID       string `json:"user_id"`
	Date         string `json:"date"`
	Time         string `json:"time"`
	Guests       int    `json:"guests"`
	Status       string `json:"status"`
}

var supabase *supabase.Client

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize Supabase client
	supabase = supabase.CreateClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"))

	// Initialize Fiber app
	app := fiber.New()

	// Define routes
	app.Get("/restaurants/:restaurantId/reservations", getReservations)
	app.Post("/restaurants/:restaurantId/reservations", createReservation)
	app.Get("/restaurants/:restaurantId/reservations/:reservationId", getReservationByID)
	app.Put("/restaurants/:restaurantId/reservations/:reservationId", updateReservation)
	app.Delete("/restaurants/:restaurantId/reservations/:reservationId", deleteReservation)

	// Start server
	port := os.Getenv("PORT")
	log.Fatal(app.Listen(":" + port))
}

// Get all reservations for a restaurant
func getReservations(c *fiber.Ctx) error {
	restaurantId := c.Params("restaurantId")
	date := c.Query("date")

	query := supabase.DB.From("reservations").Select("*").Eq("restaurant_id", restaurantId)
	if date != "" {
		query = query.Eq("date", date)
	}

	var reservations []Reservation
	if err := query.Execute(&reservations); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch reservations"})
	}

	return c.JSON(reservations)
}

// Create a new reservation
func createReservation(c *fiber.Ctx) error {
	restaurantId := c.Params("restaurantId")
	var reservation Reservation

	if err := c.BodyParser(&reservation); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	reservation.RestaurantID = restaurantId

	if err := supabase.DB.From("reservations").Insert(reservation).Execute(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create reservation"})
	}

	return c.Status(201).JSON(reservation)
}

// Get a single reservation by ID
func getReservationByID(c *fiber.Ctx) error {
	restaurantId := c.Params("restaurantId")
	reservationId := c.Params("reservationId")

	var reservation Reservation
	if err := supabase.DB.From("reservations").Select("*").Eq("restaurant_id", restaurantId).Eq("id", reservationId).Single().Execute(&reservation); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "Reservation not found"})
	}

	return c.JSON(reservation)
}

// Update a reservation
func updateReservation(c *fiber.Ctx) error {
	restaurantId := c.Params("restaurantId")
	reservationId := c.Params("reservationId")
	var updateData Reservation

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := supabase.DB.From("reservations").Update(updateData).Eq("restaurant_id", restaurantId).Eq("id", reservationId).Execute(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update reservation"})
	}

	return c.JSON(fiber.Map{"message": "Reservation updated successfully"})
}

// Delete a reservation
func deleteReservation(c *fiber.Ctx) error {
	restaurantId := c.Params("restaurantId")
	reservationId := c.Params("reservationId")

	if err := supabase.DB.From("reservations").Delete().Eq("restaurant_id", restaurantId).Eq("id", reservationId).Execute(); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete reservation"})
	}

	return c.Status(204).Send(nil)
}
