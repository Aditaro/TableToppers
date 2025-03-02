package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/supabase-community/gotrue-go/types"
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

	// Initialize the Supabase client
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

	// Create a new Gin router
	router := gin.Default()

	// Enable CORS with default settings (allow all origins)
	router.Use(cors.Default())

	// POST route to register a new user
	router.POST("/register", func(c *gin.Context) {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Create a SignupRequest struct to pass to Signup function
		signupReq := types.SignupRequest{
			Email:    request.Email,
			Password: request.Password,
		}

		// Registration using Supabase Auth
		_, err := client.Auth.Signup(signupReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	})

	// POST route to log in an existing user
	router.POST("/login", func(c *gin.Context) {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Log in the user using Supabase Auth
		session, err := client.Auth.SignInWithEmailPassword(request.Email, request.Password) // Correct method
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "session": session})
	})

	// Route for a blank homepage for now
	router.GET("/home", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the homepage!"})
	})

	// Start server
	fmt.Println("Server running on port 8080")
	router.Run(":8080")
}
