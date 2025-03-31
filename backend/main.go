package main

import (
	"bytes"
	"encoding/json"
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

type Restaurant struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	Description  string `json:"description"`
	Phone        string `json:"phone"`
	OpeningHours string `json:"opening_hours"`
	Img          string `json:"img"`
	CreatedAt    string `json:"created_at"`
}

type RestaurantCreate struct {
	Name         string `json:"name" binding:"required"`
	Location     string `json:"location" binding:"required"`
	Description  string `json:"description"`
	Phone        string `json:"phone"`
	OpeningHours string `json:"opening_hours"`
	Img          string `json:"img"`
}

type RestaurantUpdate struct {
	Name         string `json:"name"`
	Location     string `json:"location"`
	Description  string `json:"description"`
	Phone        string `json:"phone"`
	OpeningHours string `json:"opening_hours"`
	Img          string `json:"img"`
}

// Load environment variables
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

// Register Handler
func registerHandler(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		signupReq := types.SignupRequest{
			Email:    request.Email,
			Password: request.Password,
		}

		_, err := client.Auth.Signup(signupReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	}
}

// Login Handler
func loginHandler(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		session, err := client.Auth.SignInWithEmailPassword(request.Email, request.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "session": session})
	}
}

// Home Handler
func homeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the homepage!"})
	}
}

// Get Restaurants or Single Restaurant by ID Handler
func getRestaurants() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		city := c.Query("city")
		name := c.Query("name")

		url := fmt.Sprintf("%s/rest/v1/restaurants", os.Getenv("SUPABASE_URL"))
		if id != "" {
			url = fmt.Sprintf("%s/rest/v1/restaurants?id=eq.%s", os.Getenv("SUPABASE_URL"), id)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))

		if id == "" {
			query := req.URL.Query()
			if city != "" {
				query.Add("location", "eq."+city)
			}
			if name != "" {
				query.Add("name", "ilike.*"+name+"*")
			}
			req.URL.RawQuery = query.Encode()
		}

		clientHTTP := &http.Client{}
		resp, err := clientHTTP.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch restaurants"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to fetch restaurants"})
			return
		}

		var restaurants []Restaurant
		if err := json.NewDecoder(resp.Body).Decode(&restaurants); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			return
		}

		c.JSON(http.StatusOK, restaurants)
	}
}

// Create Restaurant Handler
func createRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newRestaurant RestaurantCreate

		// Validate and bind JSON input to struct
		if err := c.ShouldBindJSON(&newRestaurant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		url := fmt.Sprintf("%s/rest/v1/restaurants", os.Getenv("SUPABASE_URL"))
		requestBody, err := json.Marshal(newRestaurant)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request body"})
			return
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Content-Type", "application/json")

		clientHTTP := &http.Client{}
		resp, err := clientHTTP.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restaurant"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to create restaurant"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Restaurant created successfully"})
	}
}

// Update Restaurant Handler
func updateRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var updatedRestaurant RestaurantUpdate

		// Validate and bind JSON input to struct
		if err := c.ShouldBindJSON(&updatedRestaurant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		// Prepare request body for update
		url := fmt.Sprintf("%s/rest/v1/restaurants?id=eq.%s", os.Getenv("SUPABASE_URL"), id)
		requestBody, err := json.Marshal(updatedRestaurant)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request body: " + err.Error()})
			return
		}

		// Send PUT request to Supabase API
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request: " + err.Error()})
			return
		}

		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Content-Type", "application/json")

		log.Printf("Request Body: %s", string(requestBody))
		log.Printf("URL: %s", url)

		// Create HTTP client and send request
		clientHTTP := &http.Client{}
		resp, err := clientHTTP.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		// Check response status code
		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to update restaurant, status code: " + fmt.Sprint(resp.StatusCode)})
			return
		}

		// Return success response
		c.JSON(http.StatusOK, gin.H{"message": "Restaurant updated successfully"})
	}
}

func main() {
	loadEnv()

	client, err := initSupabase()
	if err != nil {
		log.Fatalf("Error initializing Supabase client: %v", err)
	}

	router := gin.Default()
	router.Use(cors.Default())

	router.POST("/register", registerHandler(client))
	router.POST("/login", loginHandler(client))
	router.GET("/home", homeHandler())
	router.GET("/restaurants", getRestaurants())
	router.GET("/restaurants/:id", getRestaurants())
	router.POST("/restaurants", createRestaurant())
	router.PUT("/restaurants/:id", updateRestaurant())

	fmt.Println("Server running on port 8080")
	router.Run(":8080")
}
