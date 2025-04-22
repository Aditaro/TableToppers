package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockSupabaseClient mocks Supabase interactions
type MockSupabaseClient struct {
	mock.Mock
}

// Mock method for fetching tables by restaurant ID
func (m *MockSupabaseClient) GetTablesByRestaurant(restaurantID string) ([]map[string]interface{}, error) {
	args := m.Called(restaurantID)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

// Mock method for creating a new table
func (m *MockSupabaseClient) CreateTable(restaurantID string, tableData map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(restaurantID, tableData)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// Mock method for getting a specific table
func (m *MockSupabaseClient) GetTableByID(restaurantID, tableID string) (map[string]interface{}, error) {
	args := m.Called(restaurantID, tableID)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// Mock method for updating a table
func (m *MockSupabaseClient) UpdateTable(restaurantID, tableID string, updates map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(restaurantID, tableID, updates)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// Mock method for deleting a table
func (m *MockSupabaseClient) DeleteTable(restaurantID, tableID string) error {
	args := m.Called(restaurantID, tableID)
	return args.Error(0)
}

// Mock method for fetching reservations by restaurant ID
func (m *MockSupabaseClient) GetReservationsByRestaurant(restaurantID string) ([]map[string]interface{}, error) {
	args := m.Called(restaurantID)
	return args.Get(0).([]map[string]interface{}), args.Error(1)
}

// Mock method for creating a reservation
func (m *MockSupabaseClient) CreateReservation(restaurantID string, reservationData map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(restaurantID, reservationData)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// Mock method for getting a specific reservation
func (m *MockSupabaseClient) GetReservationByID(restaurantID, reservationID string) (map[string]interface{}, error) {
	args := m.Called(restaurantID, reservationID)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// Mock method for updating a reservation
func (m *MockSupabaseClient) UpdateReservation(restaurantID, reservationID string, updates map[string]interface{}) (map[string]interface{}, error) {
	args := m.Called(restaurantID, reservationID, updates)
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

// Mock method for deleting a reservation
func (m *MockSupabaseClient) DeleteReservation(restaurantID, reservationID string) error {
	args := m.Called(restaurantID, reservationID)
	return args.Error(0)
}
