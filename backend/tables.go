package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/supabase-community/supabase-go"
)

// Table represents the structure of a restaurant table
type Table struct {
	ID           string `json:"id"`
	RestaurantID string `json:"restaurantId"`
	Number       int    `json:"number"`
	Seats        int    `json:"seats"`
	Status       string `json:"status"`
}

// Get all tables for a restaurant
func getTables(c *gin.Context) {
	restaurantID := c.Param("restaurantId")

	var tables []Table
	query := client.From("tables").Select("*").Eq("restaurant_id", restaurantID)
	if err := query.Execute(&tables); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tables)
}

// Create a new table in a restaurant
func createTable(c *gin.Context) {
	restaurantID := c.Param("restaurantId")
	var table Table

	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

table.RestaurantID = restaurantID
	_, err := client.From("tables").Insert(table, false, "", "", "").Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create table"})
		return
	}

	c.JSON(http.StatusCreated, table)
}

// Get details of a specific table
func getTableByID(c *gin.Context) {
	restaurantID := c.Param("restaurantId")
	tableID := c.Param("tableId")

	var table Table
	query := client.From("tables").Select("*").Eq("restaurant_id", restaurantID).Eq("id", tableID)
	if err := query.Execute(&table); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}

	c.JSON(http.StatusOK, table)
}

// Update a table
type TableUpdate struct {
	Number int    `json:"number,omitempty"`
	Seats  int    `json:"seats,omitempty"`
	Status string `json:"status,omitempty"`
}

func updateTable(c *gin.Context) {
	restaurantID := c.Param("restaurantId")
	tableID := c.Param("tableId")
	var updateData TableUpdate

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	_, err := client.From("tables").Update(updateData, "", "").Eq("restaurant_id", restaurantID).Eq("id", tableID).Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update table"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Table updated successfully"})
}

// Delete a table
func deleteTable(c *gin.Context) {
	restaurantID := c.Param("restaurantId")
	tableID := c.Param("tableId")

	_, err := client.From("tables").Delete("", "").Eq("restaurant_id", restaurantID).Eq("id", tableID).Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete table"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// Register table routes
func registerTableRoutes(router *gin.Engine) {
	tables := router.Group("/restaurants/:restaurantId/tables")
	{
		tables.GET("", getTables)
		tables.POST("", createTable)
		tables.GET(":tableId", getTableByID)
		tables.PUT(":tableId", updateTable)
		tables.DELETE(":tableId", deleteTable)
	}
}
