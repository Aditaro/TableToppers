package main

import (
	"tabletoppers/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTables(t *testing.T) {
	mockDB := new(mocks.MockSupabaseClient)
	mockDB.On("GetTablesByRestaurant", "restaurant1").Return([]map[string]interface{}{
		{"id": "table1", "number": 1},
		{"id": "table2", "number": 2},
	}, nil)

	tables, err := mockDB.GetTablesByRestaurant("restaurant1")

	assert.NoError(t, err)
	assert.Len(t, tables, 2)
	assert.Equal(t, "table1", tables[0]["id"])
}
