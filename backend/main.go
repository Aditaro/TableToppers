package main

import (
	"encoding/json"
	"bytes"
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
	}
}

// Login Handler
func loginHandler(client *supabase.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
	app.Post("/login", func(c *fiber.Ctx) error {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		session, err := client.Auth.SignInWithEmailPassword(request.Email, request.Password)
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

	fmt.Println("Server running on port 8080")
	router.Run(":8080")
}

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