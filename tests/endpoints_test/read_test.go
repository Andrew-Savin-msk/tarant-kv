package endpoints_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleReadKeys_Success(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}

	reqB := ReadRequest{
		Keys: []string{"key1", "key2"},
	}

	reqBody, err := json.Marshal(reqB)
	assert.NoError(t, err)
	req, err := http.NewRequest("POST", "http://localhost:8081/api/read", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	var tuple []byte
	resp.Body.Read(tuple)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, tuple)

	var respData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	assert.NoError(t, err)
}

func TestHandleReadKeys_PartialSuccess(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}
	

	reqBody := `{"keys": ["key1", "non_existing_key"]}`
	req, err := http.NewRequest("POST", "http://localhost:8081/api/read", bytes.NewBuffer([]byte(reqBody)))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var respData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	assert.NoError(t, err)
	assert.Empty(t, respData["error"], respData["error"])
	assert.NotEmpty(t, respData["not_found"])
}

func TestHandleReadKeys_BadRequest(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}

	reqBody := `{"keys": }`
	req, err := http.NewRequest("POST", "http://localhost:8081/api/read", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	var tuple []byte
	resp.Body.Read(tuple)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode, tuple)
}
