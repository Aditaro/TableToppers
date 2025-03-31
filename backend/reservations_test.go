package main

import (
	"tabletoppers/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetReservations(t *testing.T) {
	mockDB := new(mocks.MockSupabaseClient)
	mockDB.On("GetReservationsByRestaurant", "restaurant1").Return([]map[string]interface{}{
		{"id": "res1", "table_id": "table1"},
		{"id": "res2", "table_id": "table2"},
	}, nil)

	reservations, err := mockDB.GetReservationsByRestaurant("restaurant1")

	assert.NoError(t, err)
	assert.Len(t, reservations, 2)
	assert.Equal(t, "res1", reservations[0]["id"])
}
