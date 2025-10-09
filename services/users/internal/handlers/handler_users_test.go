package handlers

import (
	"encoding/json"
	"net/http"
	"service-user/internal/database"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthz(t *testing.T) {
	client := http.DefaultClient
	url := "http://localhost:8080/healthz"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(t, err)
	require.NotNil(t, req)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	type HealthzResponse struct {
		Message string `json:"message"`
	}

	var dbResp HealthzResponse
	err = json.NewDecoder(resp.Body).Decode(&dbResp)
	require.NoError(t, err)
	assert.Equal(t, "service users good", dbResp.Message)
}

func TestUsersEndpoints(t *testing.T) {
	type SignUpPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type LoginPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type SignUpResponse struct {
		User database.User `json:"user"`
	}

	type LoginResponse struct {
		User database.User `json:"user"`
	}

	client := http.DefaultClient

	// user sign up
	url := "http://localhost:8080/users"
	signUpPayload := &SignUpPayload{
		Email:    "sample@something.com",
		Password: "password123",
	}

	data, err := json.Marshal(signUpPayload)
	require.NoError(t, err)
	require.NotNil(t, data)
	require.NotEmpty(t, data)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(string(data)))
	require.NoError(t, err)
	require.NotNil(t, req)

	resp, err := client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	// defer resp.Body.Close()

	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var signupResp SignUpResponse
	err = json.NewDecoder(resp.Body).Decode(&signupResp)
	require.NoError(t, err)
	assert.Equal(t, signUpPayload.Email, signupResp.User.Email)
	resp.Body.Close()

	// user sign up, duplicate
	req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(string(data)))
	require.NoError(t, err)
	require.NotNil(t, req)

	resp, err = client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	// defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	resp.Body.Close()

	// user sign up, missing password
	signUpPayload = &SignUpPayload{
		Email:    "missingpassword@something.com",
		Password: "",
	}
	data, err = json.Marshal(signUpPayload)
	require.NoError(t, err)
	require.NotNil(t, data)
	require.NotEmpty(t, data)

	req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(string(data)))
	require.NoError(t, err)
	require.NotNil(t, req)

	resp, err = client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	// defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// user sign up, missing email
	signUpPayload = &SignUpPayload{
		Email:    "",
		Password: "missingemail",
	}
	data, err = json.Marshal(signUpPayload)
	require.NoError(t, err)
	require.NotNil(t, data)
	require.NotEmpty(t, data)

	req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(string(data)))
	require.NoError(t, err)
	require.NotNil(t, req)

	resp, err = client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	// defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	resp.Body.Close()

	// user login
	url = "http://localhost:8080/auth/login"
	loginPayload := &LoginPayload{
		Email:    "sample@something.com",
		Password: "password123",
	}

	data, err = json.Marshal(loginPayload)
	require.NoError(t, err)
	require.NotNil(t, data)
	require.NotEmpty(t, data)

	req, err = http.NewRequest(http.MethodPost, url, strings.NewReader(string(data)))
	require.NoError(t, err)
	require.NotNil(t, req)

	resp, err = client.Do(req)
	require.NoError(t, err)
	require.NotNil(t, resp)
	// defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResp LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResp)
	resp.Body.Close()

	require.NoError(t, err)
	assert.Equal(t, signupResp.User.Id, loginResp.User.Id)
	assert.Equal(t, signupResp.User.CreatedAt, loginResp.User.CreatedAt)
	assert.Equal(t, signupResp.User.UpdatedAt, loginResp.User.UpdatedAt)
	assert.Equal(t, signupResp.User.Username, loginResp.User.Username)
	assert.Equal(t, signupResp.User.Email, loginResp.User.Email)
	assert.Equal(t, signupResp.User.Role, loginResp.User.Role)
}
