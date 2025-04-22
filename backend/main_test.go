package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/gotrue-go/types"
)

type MockSupabaseClient struct {
	// Mock the methods of Supabase client
	Auth *MockAuthClient
}

type MockAuthClient struct{}

func (m *MockAuthClient) Signup(req types.SignupRequest) (*types.User, error) {
	// Simulate successful signup
	return &types.User{Email: req.Email}, nil
}

func (m *MockAuthClient) SignInWithEmailPassword(email, password string) (*types.Session, error) {
	// Simulate successful login
	if email == "test@example.com" && password == "password123" {
		return &types.Session{AccessToken: "mock-token"}, nil
	}
	return nil, fmt.Errorf("invalid credentials")
}

func TestRegisterUser(t *testing.T) {
	// Create a mock Supabase client
	client := &MockSupabaseClient{
		Auth: &MockAuthClient{},
	}

	router := gin.Default()
	router.Use(cors.Default())

	// Register route
	router.POST("/register", func(c *gin.Context) {
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

		// Simulate Supabase signup
		_, err := client.Auth.Signup(signupReq)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	})

	// Create a test request
	jsonData := `{"email": "test@example.com", "password": "password123"}`
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	expectedMessage := "User registered successfully"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %v, got %v", expectedMessage, response["message"])
	}
}

func TestLoginUser(t *testing.T) {
	// Create a mock Supabase client
	client := &MockSupabaseClient{
		Auth: &MockAuthClient{},
	}

	router := gin.Default()
	router.Use(cors.Default())

	// Login route
	router.POST("/login", func(c *gin.Context) {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Simulate Supabase login
		session, err := client.Auth.SignInWithEmailPassword(request.Email, request.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "session": session})
	})

	// Create a test request
	jsonData := `{"email": "test@example.com", "password": "password123"}`
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(jsonData)))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", w.Code)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	expectedMessage := "Login successful"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %v, got %v", expectedMessage, response["message"])
	}
}

func TestGetHomePage(t *testing.T) {
	// Set up the Gin router
	router := gin.Default()
	router.Use(cors.Default())

	// Home route
	router.GET("/home", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the homepage!"})
	})

	// Create a test request
	req, err := http.NewRequest("GET", "/home", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", w.Code)
	}

	// Check response body
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	// Validate message
	expectedMessage := "Welcome to the homepage!"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %v, got %v", expectedMessage, response["message"])
	}
}

func TestGetRestaurants(t *testing.T) {
	// Mock data for restaurants
	mockRestaurants := []Restaurant{
		{
			ID:           "1",
			Name:         "Restaurant 1",
			Location:     "City 1",
			Description:  "A nice place",
			Phone:        "123456789",
			OpeningHours: "9am - 9pm",
			Img:          "image1.png",
			CreatedAt:    "2025-03-01T12:00:00Z",
		},
		{
			ID:           "2",
			Name:         "Restaurant 2",
			Location:     "City 2",
			Description:  "Another nice place",
			Phone:        "987654321",
			OpeningHours: "10am - 10pm",
			Img:          "image2.png",
			CreatedAt:    "2025-03-01T12:00:00Z",
		},
	}

	router := gin.Default()
	router.Use(cors.Default())

	router.GET("/restaurants", func(c *gin.Context) {
		c.JSON(http.StatusOK, mockRestaurants)
	})

	req, err := http.NewRequest("GET", "/restaurants", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", w.Code)
	}

	// Check response body
	var response []Restaurant
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	assert.Equal(t, mockRestaurants, response)
}

func TestCreateRestaurant(t *testing.T) {
	// Mock data for creating a restaurant
	mockRestaurant := Restaurant{
		ID:           "3e48f81d-c6b9-4a7e-89be-c987d33c30e5",
		Name:         "Test Restaurant",
		Location:     "Test City",
		Description:  "A mock restaurant for testing",
		Phone:        "123456789",
		OpeningHours: "8am - 10pm",
		Img:          "test-image.png",
		CreatedAt:    "2025-03-31T12:00:00Z",
	}

	router := gin.Default()
	router.Use(cors.Default())

	// Mock POST handler for creating a restaurant
	router.POST("/restaurants", func(c *gin.Context) {
		var newRestaurant Restaurant
		if err := c.ShouldBindJSON(&newRestaurant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		// Simulate success and return a 201 status
		c.JSON(http.StatusCreated, gin.H{"message": "Restaurant created successfully"})
	})

	// Create a test request
	requestBody, err := json.Marshal(mockRestaurant)
	if err != nil {
		t.Fatalf("could not marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/restaurants", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %v", w.Code)
	}

	// Check response body
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	expectedMessage := "Restaurant created successfully"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %v, got %v", expectedMessage, response["message"])
	}
}

func TestDeleteRestaurant(t *testing.T) {
	// Mock data for deleting a restaurant
	restaurantID := "3e48f81d-c6b9-4a7e-89be-c987d33c30e5"

	router := gin.Default()
	router.Use(cors.Default())

	// Mock DELETE handler for deleting a restaurant
	router.DELETE("/restaurants/:id", func(c *gin.Context) {
		restaurantID := c.Param("id")
		if restaurantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id is required"})
			return
		}

		// Simulate success and return a 200 status
		c.JSON(http.StatusOK, gin.H{"message": "Restaurant deleted successfully"})
	})

	// Create a test request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/restaurants/%s", restaurantID), nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", w.Code)
	}

	// Check response body
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	expectedMessage := "Restaurant deleted successfully"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %v, got %v", expectedMessage, response["message"])
	}
}

func TestGetTables(t *testing.T) {
	// Mock data for tables
	mockTables := []Table{
		{
			ID:           "34f156e9-22d4-4507-85a2-aadd843ac251",
			RestaurantID: "059ffaf3-1409-4da1-b1c5-187dda0e27a5",
			Number:       1,
			MinCapacity:  2,
			MaxCapacity:  4,
			Status:       "available",
			X:            10,
			Y:            20,
		},
		{
			ID:           "ae5c0877-995f-43e1-8724-f81d16c38ef2",
			RestaurantID: "70336dc8-0764-432c-882c-033c2b0eac65",
			Number:       2,
			MinCapacity:  4,
			MaxCapacity:  6,
			Status:       "occupied",
			X:            15,
			Y:            25,
		},
	}

	router := gin.Default()
	router.Use(cors.Default())

	// Mock GET handler for fetching tables by restaurant ID
	router.GET("/restaurants/:id/tables", func(c *gin.Context) {
		restaurantID := c.Param("id")
		if restaurantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id is required"})
			return
		}
		// Simulate returning mock data for tables
		c.JSON(http.StatusOK, mockTables)
	})

	// Create a test request
	req, err := http.NewRequest("GET", "/restaurants/059ffaf3-1409-4da1-b1c5-187dda0e27a5/tables", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", w.Code)
	}

	// Check response body
	var response []Table
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	assert.Equal(t, mockTables, response)
}

func TestCreateTable(t *testing.T) {
	// Mock data for creating a table
	mockTable := TableCreate{
		RestaurantID: "059ffaf3-1409-4da1-b1c5-187dda0e27a5",
		Number:       1,
		MinCapacity:  2,
		MaxCapacity:  4,
		Status:       "available",
		X:            10,
		Y:            20,
	}

	router := gin.Default()
	router.Use(cors.Default())

	// Mock POST handler for creating a table
	router.POST("/restaurants/:id/tables", func(c *gin.Context) {
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

		// Simulate success and return a 201 status
		c.JSON(http.StatusCreated, gin.H{"message": "Table created successfully"})
	})

	// Create a test request
	requestBody, err := json.Marshal(mockTable)
	if err != nil {
		t.Fatalf("could not marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/restaurants/059ffaf3-1409-4da1-b1c5-187dda0e27a5/tables", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %v", w.Code)
	}

	// Check response body
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	expectedMessage := "Table created successfully"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %v, got %v", expectedMessage, response["message"])
	}
}

func TestDeleteTable(t *testing.T) {
	// Mock data for deleting a table
	tableID := "34f156e9-22d4-4507-85a2-aadd843ac251"
	restaurantID := "059ffaf3-1409-4da1-b1c5-187dda0e27a5"

	router := gin.Default()
	router.Use(cors.Default())

	// Mock DELETE handler for deleting a table
	router.DELETE("/restaurants/:id/tables/:table_id", func(c *gin.Context) {
		restaurantID := c.Param("id")
		tableID := c.Param("table_id")

		if restaurantID == "" || tableID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id and table_id are required"})
			return
		}

		// Simulate success and return a 200 status
		c.JSON(http.StatusOK, gin.H{"message": "Table deleted successfully"})
	})

	// Create a test request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/restaurants/%s/tables/%s", restaurantID, tableID), nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", w.Code)
	}

	// Check response body
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	expectedMessage := "Table deleted successfully"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %v, got %v", expectedMessage, response["message"])
	}
}

func TestGetWaitlist(t *testing.T) {
	// Mock data for waitlist entries
	mockWaitlist := []WaitlistEntry{
		{
			ID:                "34f156e9-22d4-4507-85a2-aadd843ac251",
			RestaurantID:      "059ffaf3-1409-4da1-b1c5-187dda0e27a5",
			Name:              "John Doe",
			PhoneNumber:       "1234567890",
			PartySize:         4,
			PartyAhead:        2,
			EstimatedWaitTime: 15,
			CreatedAt:         "2025-04-21T12:34:56Z",
		},
		{
			ID:                "ae5c0877-995f-43e1-8724-f81d16c38ef2",
			RestaurantID:      "70336dc8-0764-432c-882c-033c2b0eac65",
			Name:              "Jane Smith",
			PhoneNumber:       "0987654321",
			PartySize:         2,
			PartyAhead:        1,
			EstimatedWaitTime: 10,
			CreatedAt:         "2025-04-21T13:00:00Z",
		},
	}

	router := gin.Default()
	router.Use(cors.Default())

	// Mock GET handler for fetching waitlist entries by restaurant ID
	router.GET("/restaurants/:id/waitlist", func(c *gin.Context) {
		restaurantID := c.Param("id")
		if restaurantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id is required"})
			return
		}
		// Simulate returning mock data for waitlist entries
		c.JSON(http.StatusOK, mockWaitlist)
	})

	// Create a test request
	req, err := http.NewRequest("GET", "/restaurants/059ffaf3-1409-4da1-b1c5-187dda0e27a5/waitlist", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", w.Code)
	}

	// Check response body
	var response []WaitlistEntry
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	assert.Equal(t, mockWaitlist, response)
}

func TestCreateWaitlistEntry(t *testing.T) {
	// Mock data for creating a waitlist entry
	mockWaitlistEntry := WaitlistEntryCreate{
		RestaurantID:      "059ffaf3-1409-4da1-b1c5-187dda0e27a5",
		Name:              "John Doe",
		PhoneNumber:       "1234567890",
		PartySize:         4,
		PartyAhead:        2,
		EstimatedWaitTime: 15,
	}

	router := gin.Default()
	router.Use(cors.Default())

	// Mock POST handler for creating a waitlist entry
	router.POST("/restaurants/:id/waitlist", func(c *gin.Context) {
		restaurantID := c.Param("id")

		if restaurantID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id is required"})
			return
		}

		var newWaitlistEntry WaitlistEntryCreate
		if err := c.ShouldBindJSON(&newWaitlistEntry); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}

		// Simulate success and return a 201 status
		c.JSON(http.StatusCreated, gin.H{"message": "Waitlist entry created successfully"})
	})

	// Create a test request
	requestBody, err := json.Marshal(mockWaitlistEntry)
	if err != nil {
		t.Fatalf("could not marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "/restaurants/059ffaf3-1409-4da1-b1c5-187dda0e27a5/waitlist", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response status code
	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %v", w.Code)
	}

	// Check response body
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	expectedMessage := "Waitlist entry created successfully"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %v, got %v", expectedMessage, response["message"])
	}
}

func TestDeleteWaitlistEntry(t *testing.T) {
	// Mock data for deleting a waitlist entry
	waitlistID := "34f156e9-22d4-4507-85a2-aadd843ac251"
	restaurantID := "059ffaf3-1409-4da1-b1c5-187dda0e27a5"

	router := gin.Default()
	router.Use(cors.Default())

	// Mock DELETE handler for deleting a waitlist entry
	router.DELETE("/restaurants/:id/waitlist/:waitlist_id", func(c *gin.Context) {
		restaurantID := c.Param("id")
		waitlistID := c.Param("waitlist_id")

		if restaurantID == "" || waitlistID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id and waitlist_id are required"})
			return
		}

		// Simulate success and return a 200 status
		c.JSON(http.StatusOK, gin.H{"message": "Waitlist entry deleted successfully"})
	})

	// Create a test request
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/restaurants/%s/waitlist/%s", restaurantID, waitlistID), nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Record response
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Check response status code
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %v", w.Code)
	}

	// Check response body
	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("could not parse response: %v", err)
	}

	expectedMessage := "Waitlist entry deleted successfully"
	if response["message"] != expectedMessage {
		t.Errorf("expected message %v, got %v", expectedMessage, response["message"])
	}
}
