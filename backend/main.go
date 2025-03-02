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

// Restaurant details struct variable
type Restaurant struct {
	Name         string `json:"name"`
	Location     string `json:"location"`
	Description  string `json:"description,omitempty"`
	Phone        string `json:"phone,omitempty"`
	OpeningHours string `json:"openingHours,omitempty"`
	Img          string `json:"img,omitempty"` // Image URL
}


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

// Restaurant Image upload
func uploadImage(client *supabase.Client, fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	var fileBuffer bytes.Buffer
	_, err = io.Copy(&fileBuffer, file)
	if err != nil {
		return "", err
	}

	// Define file path
	filePath := "restaurants/" + fileHeader.Filename

	// Upload image to Supabase Storage
	_, err = client.Storage.From("restaurant-images").Upload(filePath, &fileBuffer, "image/jpeg")
	if err != nil {
		return "", err
	}

	// Return public URL of the uploaded image
	return client.Storage.From("restaurant-images").GetPublicURL(filePath), nil
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

		// POST route to add a restaurant
	router.POST("/restaurants", func(c *gin.Context) {
		var newRestaurant Restaurant
		if err := c.ShouldBindJSON(&newRestaurant); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// Insert data into Supabase
		query := fmt.Sprintf(`
			INSERT INTO restaurants (name, location, description, phone, opening_hours, img)
			VALUES ('%s', '%s', '%s', '%s', '%s', '%s')
			RETURNING id;
		`, newRestaurant.Name, newRestaurant.Location, newRestaurant.Description,
		newRestaurant.Phone, newRestaurant.OpeningHours, newRestaurant.Img)

		var restaurantID string
		err := client.DB.From("restaurants").Execute(query).Scan(&restaurantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newRestaurant.ID = restaurantID
		c.JSON(http.StatusOK, gin.H{"message": "Restaurant added successfully", "restaurant": newRestaurant})
	})

	// Route for adding restaurants
	router.POST("/restaurants", func(c *gin.Context) {
		var restaurant Restaurant
	
		// Parse form data
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "File too large"})
			return
		}
	
		// Get restaurant details
		restaurant.Name = c.PostForm("name")
		restaurant.Location = c.PostForm("location")
		restaurant.Description = c.PostForm("description")
		restaurant.Phone = c.PostForm("phone")
		restaurant.OpeningHours = c.PostForm("openingHours")
	
		// Handle image upload
		file, fileHeader, err := c.Request.FormFile("img")
		if err == nil { // Image is optional
			defer file.Close()
			imgURL, uploadErr := uploadImage(client, fileHeader)
			if uploadErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Image upload failed"})
				return
			}
			restaurant.Img = imgURL
		}
	
		// Save restaurant data to Supabase
		_, err = client.From("restaurants").Insert(restaurant, false, "", "", "").Execute()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save restaurant"})
			return
		}
	
		c.JSON(http.StatusOK, gin.H{"message": "Restaurant added successfully", "restaurant": restaurant})
	})	

	// Route for a blank homepage for now
	router.GET("/home", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the homepage!"})
	})

	// Start server
	fmt.Println("Server running on port 8080")
	router.Run(":8080")
}
