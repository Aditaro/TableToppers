package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
	"github.com/supabase-community/gotrue-go/types"
)

// Restaurant details struct variable
type Restaurant struct {
	Name         string `json:"name"`
	Location     string `json:"location"`
	Description  string `json:"description,omitempty"`
	Phone        string `json:"phone,omitempty"`
	OpeningHours string `json:"openingHours,omitempty"`
	Img          string `json:"img,omitempty"` // Image URL
}

// Table struct updated to match actual database schema
type Table struct {
	ID           string `json:"id"`
	RestaurantID string `json:"restaurant_id"`
	Number       int    `json:"number"`
	MinCapacity  int    `json:"min_capacity"`
	MaxCapacity  int    `json:"max_capacity"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
	Status       string `json:"status"`
}

// Load environment variables from a .env file
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

// Initialize Supabase client
func initSupabase() (*supabase.Client, error) {
	url := os.Getenv("SUPABASE_URL")
	anonKey := os.Getenv("SUPABASE_ANON_KEY")

	client, err := supabase.NewClient(url, anonKey, &supabase.ClientOptions{})
	if err != nil {
		return nil, fmt.Errorf("error creating Supabase client: %v", err)
	}

	return client, nil
}

func main() {
	// Load environment variables
	loadEnv()

	// Initialize Supabase client
	client, err := initSupabase()
	if err != nil {
		log.Fatalf("Error initializing Supabase client: %v", err)
	}

	// Initialize the Fiber app
	app := fiber.New()

	// Restaurant management routes
	app.Post("/register", func(c *fiber.Ctx) error {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		signupReq := types.SignupRequest{
			Email:    request.Email,
			Password: request.Password,
		}

		_, err := client.Auth.Signup(signupReq)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User registered successfully"})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		session, err := client.Auth.SignInWithEmailPassword(request.Email, request.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login successful", "session": session})
	})

	// Restaurant creation route
	app.Post("/restaurants", func(c *fiber.Ctx) error {
		var restaurant Restaurant

		if err := c.BodyParser(&restaurant); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid form submission"})
		}

		// Insert restaurant into Supabase
		data, count, err := client.From("restaurants").Insert(restaurant, false, "", "", "").Execute()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to save restaurant: %v", err)})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Restaurant added successfully", "data": data, "count": count})
	})

	// Table routes

	// Get tables for a specific restaurant
	app.Get("/restaurants/:restaurantId/tables", func(c *fiber.Ctx) error {
		restaurantID := c.Params("restaurantId")

		var tables []Table
		// Correct syntax for Select: Select(columns string, head string, count bool)
		data, count, err := client.From("tables").Select("*", "", false).Eq("restaurant_id", restaurantID).Execute()
		if err != nil {
			fmt.Printf("Error fetching tables: %v\n", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch tables: " + err.Error()})
		}

		// Log the raw response for debugging
		fmt.Printf("Raw table data response: %s\n", string(data))

		// Unmarshal the response data
		if err := json.Unmarshal(data, &tables); err != nil {
			fmt.Printf("Error unmarshaling tables data: %v\n", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to parse response: " + err.Error()})
		}

		return c.JSON(fiber.Map{"tables": tables, "count": count})
	})

	// Create a new table for a restaurant
	app.Post("/restaurants/:restaurantId/tables", func(c *fiber.Ctx) error {
		restaurantID := c.Params("restaurantId")
		var table Table

		if err := c.BodyParser(&table); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
		}

		table.RestaurantID = restaurantID

		// Log the table data being sent
		tableJSON, _ := json.Marshal(table)
		fmt.Printf("Creating table with data: %s\n", string(tableJSON))

		// Correct syntax for Insert
		data, count, err := client.From("tables").Insert(table, false, "", "", "").Execute()
		if err != nil {
			fmt.Printf("Error creating table: %v\n", err)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create table: " + err.Error()})
		}

		return c.Status(201).JSON(fiber.Map{"message": "Table created successfully", "data": data, "count": count})
	})

	// Update an existing table for a restaurant
	app.Put("/restaurants/:restaurantId/tables/:tableId", func(c *fiber.Ctx) error {
		tableId := c.Params("tableId")
		restaurantId := c.Params("restaurantId")
		
		fmt.Printf("Updating table ID: %s for restaurant ID: %s\n", tableId, restaurantId)
		
		var updatedTable Table

		if err := c.BodyParser(&updatedTable); err != nil {
			fmt.Printf("Error parsing request body: %v\n", err)
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request body: " + err.Error()})
		}

		// Log the table data being sent for debugging
		tableJSON, _ := json.Marshal(updatedTable)
		fmt.Printf("Updating table with data: %s\n", string(tableJSON))

		// First check if we can retrieve the table
		getResult, _, err := client.From("tables").Select("*", "", false).Eq("id", tableId).Execute()
		if err != nil {
			fmt.Printf("Failed to retrieve table: %v\n", err)
			return c.Status(404).JSON(fiber.Map{"error": "Table not found: " + err.Error()})
		}
		fmt.Printf("Found table: %s\n", string(getResult))

		// Try updating with just the status field first
		updateData := map[string]string{"status": updatedTable.Status}
		updateJSON, _ := json.Marshal(updateData)
		fmt.Printf("Trying minimal update with: %s\n", string(updateJSON))
		
		data, count, err := client.From("tables").Update(updateData, "", "").Eq("id", tableId).Execute()
		if err != nil {
			fmt.Printf("Supabase update error (minimal): %v\n", err)
			
			// Try with full table data as fallback
			data, count, err = client.From("tables").Update(updatedTable, "", "").Eq("id", tableId).Execute()
			if err != nil {
				fmt.Printf("Supabase update error (full): %v\n", err)
				return c.Status(500).JSON(fiber.Map{"error": "Failed to update table: " + err.Error()})
			}
		}

		return c.Status(200).JSON(fiber.Map{"message": "Table updated successfully", "data": data, "count": count})
	})

	// Start the Fiber app
	fmt.Println("Server running on port 8082")
	if err := app.Listen(":8083"); err != nil {
		log.Fatal("Error running the server: ", err)
	}
}