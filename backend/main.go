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

// Restaurant struct for database operations
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

// RestaurantCreate struct for creation requests
type RestaurantCreate struct {
	Name         string `json:"name" binding:"required"`
	Location     string `json:"location" binding:"required"`
	Description  string `json:"description"`
	Phone        string `json:"phone"`
	OpeningHours string `json:"opening_hours"`
	Img          string `json:"img"`
}

// RestaurantUpdate struct for update requests
type RestaurantUpdate struct {
	Name         string `json:"name"`
	Location     string `json:"location"`
	Description  string `json:"description"`
	Phone        string `json:"phone"`
	OpeningHours string `json:"opening_hours"`
	Img          string `json:"img"`
}

// Table struct for database operations
type Table struct {
	ID           string `json:"id"`
	RestaurantID string `json:"restaurant_id"`
	Number       int    `json:"number"`
	MinCapacity  int    `json:"min_capacity"`
	MaxCapacity  int    `json:"max_capacity"`
	Status       string `json:"status"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
}

// TableCreate struct for creation requests
type TableCreate struct {
	RestaurantID string `json:"restaurant_id"`
	Number       int    `json:"number" binding:"required"`
	MinCapacity  int    `json:"min_capacity" binding:"required"`
	MaxCapacity  int    `json:"max_capacity" binding:"required"`
	Status       string `json:"status" binding:"required"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
}

// TableUpdate struct for update requests
type TableUpdate struct {
	Number      int    `json:"number"`
	MinCapacity int    `json:"min_capacity"`
	MaxCapacity int    `json:"max_capacity"`
	Status      string `json:"status"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
}

// Load environment variables
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Error loading .env file, using environment variables instead")
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
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the restaurant management API!"})
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
		req.Header.Set("Prefer", "return=representation")

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

		if err := c.ShouldBindJSON(&updatedRestaurant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		url := fmt.Sprintf("%s/rest/v1/restaurants?id=eq.%s", os.Getenv("SUPABASE_URL"), id)
		requestBody, err := json.Marshal(updatedRestaurant)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request body: " + err.Error()})
			return
		}

		req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request: " + err.Error()})
			return
		}

		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Prefer", "return=representation")

		log.Printf("Request Body: %s", string(requestBody))
		log.Printf("URL: %s", url)

		clientHTTP := &http.Client{}
		resp, err := clientHTTP.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to update restaurant, status code: " + fmt.Sprint(resp.StatusCode)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Restaurant updated successfully"})
	}
}

// Delete Restaurant Handler
func deleteRestaurant() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		url := fmt.Sprintf("%s/rest/v1/restaurants?id=eq.%s", os.Getenv("SUPABASE_URL"), id)

		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request: " + err.Error()})
			return
		}

		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Content-Type", "application/json")

		clientHTTP := &http.Client{}
		resp, err := clientHTTP.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request: " + err.Error()})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to delete restaurant, status code: " + fmt.Sprint(resp.StatusCode)})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Restaurant deleted successfully"})
	}
}

// Get Tables for a Specific Restaurant Handler
func getTables() gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantID := c.Param("id")

		if restaurantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id is required"})
			return
		}

		url := fmt.Sprintf("%s/rest/v1/tables?restaurant_id=eq.%s", os.Getenv("SUPABASE_URL"), restaurantID)

		// If table_id is provided, add it to the query
		if tableID := c.Param("table_id"); tableID != "" {
			url = fmt.Sprintf("%s&id=eq.%s", url, tableID)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))

		clientHTTP := &http.Client{}
		resp, err := clientHTTP.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tables"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to fetch tables"})
			return
		}

		var tables []Table
		if err := json.NewDecoder(resp.Body).Decode(&tables); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response"})
			return
		}

		c.JSON(http.StatusOK, tables)
	}
}

// Create Table Handler
func createTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantID := c.Param("id")

		if restaurantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id is required"})
			return
		}

		var newTable TableCreate
		if err := c.ShouldBindJSON(&newTable); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}
		newTable.RestaurantID = restaurantID

		url := fmt.Sprintf("%s/rest/v1/tables", os.Getenv("SUPABASE_URL"))
		requestBody, _ := json.Marshal(newTable)

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Request creation failed"})
			return
		}
		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Prefer", "return=representation")

		clientHTTP := &http.Client{}
		resp, err := clientHTTP.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create table"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to create table"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Table created successfully"})
	}
}


// Update Table Handler
func updateTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantID := c.Param("id")
		tableID := c.Param("table_id")
		var updatedTable TableUpdate

		if err := c.ShouldBindJSON(&updatedTable); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		url := fmt.Sprintf("%s/rest/v1/tables?id=eq.%s&restaurant_id=eq.%s", os.Getenv("SUPABASE_URL"), tableID, restaurantID)
		requestBody, err := json.Marshal(updatedTable)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request body"})
			return
		}

		req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Prefer", "return=representation")

		clientHTTP := &http.Client{}
		resp, err := clientHTTP.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to update table"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Table updated successfully"})
	}
}

// Delete Table Handler
func deleteTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract restaurant_id and table_id from the URL parameters
		restaurantID := c.Param("id")
		tableID := c.Param("table_id")

		// Validate restaurant_id and table_id
		if restaurantID == "" || tableID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id and table_id are required"})
			return
		}

		// Build URL to delete the table for the specific restaurant
		url := fmt.Sprintf("%s/rest/v1/tables?restaurant_id=eq.%s&id=eq.%s", os.Getenv("SUPABASE_URL"), restaurantID, tableID)

		// Create the DELETE request
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))

		clientHTTP := &http.Client{}
		resp, err := clientHTTP.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete table"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to delete table"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Table deleted successfully"})
	}
}

func main() {
	// Load environment variables
	loadEnv()

	// Initialize Supabase client
	client, err := initSupabase()
	if err != nil {
		log.Fatalf("Error initializing Supabase client: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()
	router.Use(cors.Default())

	// User routes
	router.POST("/register", registerHandler(client))
	router.POST("/login", loginHandler(client))
	router.GET("/home", homeHandler())

	// Restaurant routes
	router.GET("/restaurants", getRestaurants())
	router.GET("/restaurants/:id", getRestaurants())
	router.POST("/restaurants", createRestaurant())
	router.PUT("/restaurants/:id", updateRestaurant())
	router.DELETE("/restaurants/:id", deleteRestaurant())

	// Table routes
	router.GET("/restaurants/:id/tables", getTables())
	router.GET("/restaurants/:id/tables/:table_id", getTables())
	router.POST("/restaurants/:id/tables", createTable())
	router.PUT("/restaurants/:id/tables/:table_id", updateTable())
	router.DELETE("/restaurants/:id/tables/:table_id", deleteTable())

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8087"
	}
	
	fmt.Printf("Server running on port %s\n", port)
	router.Run(":" + port)
}