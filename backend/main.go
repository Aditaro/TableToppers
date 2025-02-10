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

// Restaurant struct to store restaurant details
type Restaurant struct {
	ID          string `json:"id"`          // Unique restaurant ID
	Name        string `json:"name"`        // Name of the restaurant
	OwnerID     string `json:"owner_id"`    // Owner's ID
	TotalTables int    `json:"total_tables"`// Number of tables available
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

	// GET route to fetch the list of restaurants
	router.GET("/restaurants", func(c *gin.Context) {
		var restaurants []Restaurant

		// Fetch all restaurants from Supabase
		err := client.From("restaurants").Select("*").Execute(&restaurants)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch restaurants"})
			return
		}

		// Return the list of restaurants
		c.JSON(http.StatusOK, gin.H{"restaurants": restaurants})
	})

	// Start server
	fmt.Println("Server running on port 8080")
	router.Run(":8080")
}
