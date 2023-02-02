package tests

import (
	"bytes"
	"encoding/json"
	"example/account-management/models"
	"example/account-management/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	router := router.SetupRouter()

	w := httptest.NewRecorder()

	body := []byte(`{"username": "timmy", "password": "test"}`)
	bodyReader := bytes.NewReader(body)
	req, _ := http.NewRequest("POST", "/users", bodyReader)
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	// Convert the JSON response to a map
	var response models.UserResponse
	json.Unmarshal([]byte(w.Body.String()), &response)

	// Grab the value & whether or not it exists
	username := response.Username
	assert.Equal(t, "timmy", username)

	score := response.Score
	assert.Equal(t, 0, score)
}
