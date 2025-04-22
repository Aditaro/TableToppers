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

// Update Restaurant Handler (DOES NOT WORK)
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

type Table struct {
	ID           string `json:"id"`
	RestaurantID string `json:"restaurant_id"`
	Number       int8   `json:"number"`
	MinCapacity  int    `json:"min_capacity"`
	MaxCapacity  int    `json:"max_capacity"`
	Status       string `json:"status"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
}

type TableCreate struct {
	RestaurantID string `json:"restaurant_id"`
	Number       int8   `json:"number" binding:"required"`
	MinCapacity  int    `json:"min_capacity" binding:"required"`
	MaxCapacity  int    `json:"max_capacity" binding:"required"`
	Status       string `json:"status" binding:"required"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
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

		requestBody, err := json.Marshal(newTable)
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

// Delete table handler
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

/*
	func updateTable() gin.HandlerFunc {
		return func(c *gin.Context) {
			restaurantID := c.Param("restaurantId")
			tableID := c.Param("tableId")
			var updatedTable TableUpdate

			if err := c.ShouldBindJSON(&updatedTable); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
				return
			}

			url := fmt.Sprintf("%s/rest/v1/tables?id=eq.%s&restaurantId=eq.%s", os.Getenv("SUPABASE_URL"), tableID, restaurantID)
			requestBody, err := json.Marshal(updatedTable)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request body"})
				return
			}

			req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
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
*/

type WaitlistEntry struct {
	ID                string `json:"id"`
	RestaurantID      string `json:"restaurant_id"`
	Name              string `json:"name"`
	PhoneNumber       string `json:"phone_number"`
	PartySize         int    `json:"party_size"`
	PartyAhead        int    `json:"party_ahead"`
	EstimatedWaitTime int    `json:"estimated_wait_time"`
	CreatedAt         string `json:"created_at"`
}

// WaitlistEntryCreate
type WaitlistEntryCreate struct {
	RestaurantID      string `json:"restaurant_id"`       // The restaurant the waitlist entry belongs to
	Name              string `json:"name"`                // Name of the person on the waitlist
	PhoneNumber       string `json:"phone_number"`        // Contact number of the person
	PartySize         int    `json:"party_size"`          // The size of the party
	PartyAhead        int    `json:"party_ahead"`         // Number of parties ahead in the waitlist
	EstimatedWaitTime int    `json:"estimated_wait_time"` // Estimated wait time in minutes
}

// Get waitlist entries for a specific restaurant Handler
func getWaitlist() gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantID := c.Param("id")

		if restaurantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id is required"})
			return
		}

		url := fmt.Sprintf("%s/rest/v1/waitlist?restaurant_id=eq.%s", os.Getenv("SUPABASE_URL"), restaurantID)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
			return
		}

		req.Header.Set("apikey", os.Getenv("SUPABASE_ANON_KEY"))
		req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_ANON_KEY"))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch waitlist entries"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to fetch waitlist entries"})
			return
		}

		var entries []WaitlistEntry
		if err := json.NewDecoder(resp.Body).Decode(&entries); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse response: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, entries)
	}
}

// Create waitlist entry for a specific restaurant handler
func createWaitlistEntry() gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantID := c.Param("id")

		if restaurantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id is required"})
			return
		}

		var newEntry WaitlistEntryCreate
		if err := c.ShouldBindJSON(&newEntry); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		newEntry.RestaurantID = restaurantID

		url := fmt.Sprintf("%s/rest/v1/waitlist", os.Getenv("SUPABASE_URL"))

		requestBody, err := json.Marshal(newEntry)
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create waitlist entry"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to create waitlist entry"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"message": "Waitlist entry created successfully"})
	}
}

// Delete waitlist entry for a specific restaurant Handler
func deleteWaitlistEntry() gin.HandlerFunc {
	return func(c *gin.Context) {
		restaurantID := c.Param("id")
		entryID := c.Param("entry_id")

		if restaurantID == "" || entryID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id and entry_id are required"})
			return
		}

		url := fmt.Sprintf("%s/rest/v1/waitlist?restaurant_id=eq.%s&id=eq.%s", os.Getenv("SUPABASE_URL"), restaurantID, entryID)

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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete waitlist entry"})
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			c.JSON(resp.StatusCode, gin.H{"error": "Failed to delete waitlist entry"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Waitlist entry deleted successfully"})
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
	router.POST("/restaurants/:id/tables", createTable())
	router.DELETE("/restaurants/:id/tables/:table_id", deleteTable())
	router.GET("/restaurants/:id/tables/:table_id", getTables())

	/*
		router.GET("/restaurants/:id/tables/:tableId", getTables())      // Get details of a specific table
		router.PUT("/restaurants/:id/tables/:tableId", updateTable())    // Update a specific table
	*/

	// Waitlist routes
	router.GET("/restaurants/:id/waitlist", getWaitlist())
	router.POST("/restaurants/:id/waitlist", createWaitlistEntry())
	router.DELETE("/restaurants/:id/waitlist/:entry_id", deleteWaitlistEntry())

	fmt.Println("Server running on port 8080")
	router.Run(":8080")
}
