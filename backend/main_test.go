package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid" // Import the CORRECT uuid package
	"github.com/stretchr/testify/assert"
	"github.com/supabase-community/gotrue-go/types"
)

// --- Mock Supabase Client ---
type MockSupabaseClient struct {
	Auth *MockAuthClient
}

type MockAuthClient struct{}

// Mock Signup method - CORRECTED with google/uuid
func (m *MockAuthClient) Signup(req types.SignupRequest) (*types.User, error) {
	if req.Email == "test@example.com" {
		// Generate a UUID using google/uuid
		mockID := uuid.New() // Generate a new random UUID (google/uuid version)

		return &types.User{
			ID:    mockID, // Assign the generated google/uuid.UUID
			Email: req.Email,
			// Add other necessary fields from types.User if needed by your tests
		}, nil
	}
	return nil, fmt.Errorf("mock signup error: invalid email")
}

// Mock SignInWithEmailPassword method - CORRECTED with google/uuid
func (m *MockAuthClient) SignInWithEmailPassword(email, password string) (*types.Session, error) {
	if email == "test@example.com" && password == "password123" {
		// Generate a UUID using google/uuid
		mockUserID := uuid.New() // Generate a new random UUID (google/uuid version)

		// Create the User struct VALUE
		mockUser := types.User{
			ID:    mockUserID, // Assign the generated google/uuid.UUID
			Email: email,
			// Add other necessary fields from types.User if needed
		}

		// Return the Session pointer, assigning the User struct VALUE
		return &types.Session{
			AccessToken: "mock-access-token",
			User:        mockUser, // Assign the value, not a pointer (&mockUser)
			// Add other necessary fields from types.Session if needed
		}, nil
	}
	return nil, fmt.Errorf("invalid credentials")
}

// --- Test Setup ---
func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// --- Test Functions ---

// TestHomeHandler remains the same
func TestHomeHandler(t *testing.T) {
	// Arrange
	router := setupRouter()
	router.GET("/home", homeHandler()) // Use actual handler

	req, _ := http.NewRequest(http.MethodGet, "/home", nil)
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	expectedBody := gin.H{"message": "Welcome to the restaurant management API!"}
	var actualBody gin.H
	err := json.Unmarshal(rr.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedBody, actualBody)
}

// TestRegisterHandler_Success remains the same (uses the corrected mock)
func TestRegisterHandler_Success(t *testing.T) {
	// Arrange
	router := setupRouter()
	mockSupabase := &MockSupabaseClient{Auth: &MockAuthClient{}}

	// Using inline handler with mock
	router.POST("/register", func(c *gin.Context) {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		signupReq := types.SignupRequest{Email: request.Email, Password: request.Password}
		_, err := mockSupabase.Auth.Signup(signupReq) // Use corrected mock
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
	})

	registerData := gin.H{"email": "test@example.com", "password": "password123"}
	bodyBytes, _ := json.Marshal(registerData)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	expectedBody := gin.H{"message": "User registered successfully"}
	var actualBody gin.H
	err := json.Unmarshal(rr.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedBody, actualBody)
}

// TestLoginHandler_Success remains the same (uses the corrected mock)
func TestLoginHandler_Success(t *testing.T) {
	// Arrange
	router := setupRouter()
	mockSupabase := &MockSupabaseClient{Auth: &MockAuthClient{}}

	// Using inline handler with mock
	router.POST("/login", func(c *gin.Context) {
		var request struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}
		session, err := mockSupabase.Auth.SignInWithEmailPassword(request.Email, request.Password) // Use corrected mock
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Login successful", "session": session})
	})

	loginData := gin.H{"email": "test@example.com", "password": "password123"}
	bodyBytes, _ := json.Marshal(loginData)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	// Act
	router.ServeHTTP(rr, req)

	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	var actualBody map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Equal(t, "Login successful", actualBody["message"])
	assert.NotNil(t, actualBody["session"])
	sessionMap, ok := actualBody["session"].(map[string]interface{})
	assert.True(t, ok, "Session should be a map")
	assert.NotEmpty(t, sessionMap["access_token"])
	assert.NotNil(t, sessionMap["user"])
	userMap, ok := sessionMap["user"].(map[string]interface{})
	assert.True(t, ok, "User in session should be a map")
	// Check the ID exists and is a string (UUIDs marshal to strings in JSON)
	userID, ok := userMap["id"].(string)
	assert.True(t, ok, "User ID should be a string in JSON")
	assert.NotEmpty(t, userID, "User ID should not be empty")
}


// --- Tests Mocking Handler Responses (No external calls) ---
// (These tests remain unchanged)

func TestGetRestaurants_Success(t *testing.T) {
	// Arrange
	router := setupRouter()
	mockRestaurants := []Restaurant{
		{ID: "res-1", Name: "Mock Cafe", Location: "Testville"},
		{ID: "res-2", Name: "Fake Bistro", Location: "Testville"},
	}
	router.GET("/restaurants", func(c *gin.Context) {
		c.JSON(http.StatusOK, mockRestaurants)
	})
	req, _ := http.NewRequest(http.MethodGet, "/restaurants", nil)
	rr := httptest.NewRecorder()
	// Act
	router.ServeHTTP(rr, req)
	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	var actualBody []Restaurant
	err := json.Unmarshal(rr.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Equal(t, mockRestaurants, actualBody)
}

func TestGetSingleRestaurant_Success(t *testing.T) {
	// Arrange
	router := setupRouter()
	mockRestaurant := Restaurant{ID: "res-1", Name: "Mock Cafe", Location: "Testville"}
	restaurantID := "res-1"
	router.GET("/restaurants/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == restaurantID {
			c.JSON(http.StatusOK, []Restaurant{mockRestaurant})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		}
	})
	req, _ := http.NewRequest(http.MethodGet, "/restaurants/"+restaurantID, nil)
	rr := httptest.NewRecorder()
	// Act
	router.ServeHTTP(rr, req)
	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	var actualBody []Restaurant
	err := json.Unmarshal(rr.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Len(t, actualBody, 1)
	assert.Equal(t, mockRestaurant, actualBody[0])
}

func TestCreateRestaurant_Success(t *testing.T) {
	// Arrange
	router := setupRouter()
	router.POST("/restaurants", func(c *gin.Context) {
		var createReq RestaurantCreate
		if err := c.ShouldBindJSON(&createReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}
		if createReq.Name == "" || createReq.Location == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Restaurant created successfully"})
	})
	createData := RestaurantCreate{Name: "New Mock Grill", Location: "Mock City"}
	bodyBytes, _ := json.Marshal(createData)
	req, _ := http.NewRequest(http.MethodPost, "/restaurants", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	// Act
	router.ServeHTTP(rr, req)
	// Assert
	assert.Equal(t, http.StatusCreated, rr.Code)
	expectedBody := gin.H{"message": "Restaurant created successfully"}
	var actualBody gin.H
	err := json.Unmarshal(rr.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedBody, actualBody)
}

func TestUpdateRestaurant_Success(t *testing.T) {
	// Arrange
	router := setupRouter()
	restaurantID := "res-to-update-123"
	router.PATCH("/restaurants/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id != restaurantID {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		var updateReq RestaurantUpdate
		if err := c.ShouldBindJSON(&updateReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Restaurant updated successfully"})
	})
	updateData := RestaurantUpdate{Description: "Updated Description"}
	bodyBytes, _ := json.Marshal(updateData)
	req, _ := http.NewRequest(http.MethodPatch, "/restaurants/"+restaurantID, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	// Act
	router.ServeHTTP(rr, req)
	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	expectedBody := gin.H{"message": "Restaurant updated successfully"}
	var actualBody gin.H
	err := json.Unmarshal(rr.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedBody, actualBody)
}

func TestDeleteRestaurant_Success(t *testing.T) {
	// Arrange
	router := setupRouter()
	restaurantID := "res-to-delete-456"
	router.DELETE("/restaurants/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == restaurantID {
			c.JSON(http.StatusOK, gin.H{"message": "Restaurant deleted successfully"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		}
	})
	req, _ := http.NewRequest(http.MethodDelete, "/restaurants/"+restaurantID, nil)
	rr := httptest.NewRecorder()
	// Act
	router.ServeHTTP(rr, req)
	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	expectedBody := gin.H{"message": "Restaurant deleted successfully"}
	var actualBody gin.H
	err := json.Unmarshal(rr.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedBody, actualBody)
}

func TestGetTables_Success(t *testing.T) {
	// Arrange
	router := setupRouter()
	restaurantID := "res-1"
	mockTables := []Table{
		{ID: "tbl-1", RestaurantID: restaurantID, Number: 1, Status: "available"},
		{ID: "tbl-2", RestaurantID: restaurantID, Number: 2, Status: "occupied"},
	}
	router.GET("/restaurants/:id/tables", func(c *gin.Context) {
		id := c.Param("id")
		if id == restaurantID {
			c.JSON(http.StatusOK, mockTables)
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
		}
	})
	req, _ := http.NewRequest(http.MethodGet, "/restaurants/"+restaurantID+"/tables", nil)
	rr := httptest.NewRecorder()
	// Act
	router.ServeHTTP(rr, req)
	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	var actualBody []Table
	err := json.Unmarshal(rr.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Equal(t, mockTables, actualBody)
}

func TestCreateTable_Success(t *testing.T) {
	// Arrange
	router := setupRouter()
	restaurantID := "res-1"
	router.POST("/restaurants/:id/tables", func(c *gin.Context) {
		id := c.Param("id")
		if id != restaurantID {
			c.JSON(http.StatusNotFound, gin.H{"error": "Restaurant not found"})
			return
		}
		var createReq TableCreate
		if err := c.ShouldBindJSON(&createReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if createReq.Number <= 0 || createReq.MinCapacity <= 0 || createReq.MaxCapacity < createReq.MinCapacity || createReq.Status == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table data"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"message": "Table created successfully"})
	})
	createData := TableCreate{Number: 5, MinCapacity: 2, MaxCapacity: 4, Status: "available"}
	bodyBytes, _ := json.Marshal(createData)
	req, _ := http.NewRequest(http.MethodPost, "/restaurants/"+restaurantID+"/tables", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	// Act
	router.ServeHTTP(rr, req)
	// Assert
	assert.Equal(t, http.StatusCreated, rr.Code)
	expectedBody := gin.H{"message": "Table created successfully"}
	var actualBody gin.H
	err := json.Unmarshal(rr.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedBody, actualBody)
}

func TestUpdateTable_Success(t *testing.T) {
	// Arrange
	router := setupRouter()
	restaurantID := "res-1"
	tableID := "tbl-to-update-abc"
	router.PUT("/restaurants/:id/tables/:table_id", func(c *gin.Context) {
		resID := c.Param("id")
		tblID := c.Param("table_id")
		if resID != restaurantID || tblID != tableID {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
			return
		}
		var updateReq TableUpdate
		if err := c.ShouldBindJSON(&updateReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Table updated successfully"})
	})
	updateData := TableUpdate{Status: "reserved", X: 50, Y: 50}
	bodyBytes, _ := json.Marshal(updateData)
	req, _ := http.NewRequest(http.MethodPut, "/restaurants/"+restaurantID+"/tables/"+tableID, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	// Act
	router.ServeHTTP(rr, req)
	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	expectedBody := gin.H{"message": "Table updated successfully"}
	var actualBody gin.H
	err := json.Unmarshal(rr.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedBody, actualBody)
}

func TestDeleteTable_Success(t *testing.T) {
	// Arrange
	router := setupRouter()
	restaurantID := "res-1"
	tableID := "tbl-to-delete-xyz"
	router.DELETE("/restaurants/:id/tables/:table_id", func(c *gin.Context) {
		resID := c.Param("id")
		tblID := c.Param("table_id")
		if resID == restaurantID && tblID == tableID {
			c.JSON(http.StatusOK, gin.H{"message": "Table deleted successfully"})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		}
	})
	req, _ := http.NewRequest(http.MethodDelete, "/restaurants/"+restaurantID+"/tables/"+tableID, nil)
	rr := httptest.NewRecorder()
	// Act
	router.ServeHTTP(rr, req)
	// Assert
	assert.Equal(t, http.StatusOK, rr.Code)
	expectedBody := gin.H{"message": "Table deleted successfully"}
	var actualBody gin.H
	err := json.Unmarshal(rr.Body.Bytes(), &actualBody)
	assert.NoError(t, err)
	assert.Equal(t, expectedBody, actualBody)
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
=======
>>>>>>> main
