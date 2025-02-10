package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

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

// Reservation struct for table booking
type Reservation struct {
	CustomerID string `json:"customer_id"`
	Date       string `json:"date"`
	Day        string `json:"day"`
	Time       string `json:"time"`
}

func main() {
	// Load environment variables
	loadEnv()

	// Initialize Supabase client
	client, err := initSupabase()
	if err != nil {
		log.Fatalf("Error initializing Supabase client: %v", err)
	}

	// Create a new Gin router
	router := gin.Default()

	// Enable CORS
	router.Use(cors.Default())

	// POST route to register a new user
	router.POST("/register", func(c *gin.Context) {
		// Registration logic...
	})

	// POST route to log in an existing user
	router.POST("/login", func(c *gin.Context) {
		// Login logic...
	})

	// POST route to reserve a table
	router.POST("/reserve", func(c *gin.Context) {
		var reservation Reservation

		// Parse incoming request
		if err := c.ShouldBindJSON(&reservation); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		// Validate reservation fields
		if reservation.Date == "" || reservation.Day == "" || reservation.Time == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Date, day, and time are required"})
			return
		}

		// Example: Simulate saving to Supabase (use actual database insert in production)
		_, err := client.From("reservations").Insert([]interface{}{
			map[string]interface{}{
				"customer_id": reservation.CustomerID,
				"date":        reservation.Date,
				"day":         reservation.Day,
				"time":        reservation.Time,
			},
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save reservation"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Reservation successful", "reservation": reservation})
	})

	// Route for a blank homepage
	router.GET("/home", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the homepage!"})
	})

	// Start server
	fmt.Println("Server running on port 8080")
	router.Run(":8080")
}
