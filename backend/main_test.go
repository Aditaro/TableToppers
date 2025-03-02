package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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
