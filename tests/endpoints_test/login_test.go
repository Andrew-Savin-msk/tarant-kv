package endpoints_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleLogin_Success(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}

	reqBody := `{"username":"admin", "password":"presale"}`
	req, err := http.NewRequest("POST", "http://localhost:8081/api/login", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	assert.NoError(t, err)
	assert.NotEmpty(t, respData["token"])
}

func TestHandleLogin_InvalidCredentials(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}

	reqBody := `{"username":"wronguser", "password":"wrongpassword"}`
	req, err := http.NewRequest("POST", "http://localhost:8081/api/login", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestHandleLogin_MissingCredentials(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}

	reqBody := `{"username":"", "password":""}`
	req, err := http.NewRequest("POST", "http://localhost:8081/api/login", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
