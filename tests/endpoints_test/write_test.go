package endpoints_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHandleWriteKeys_Success(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}

	reqB := WriteReq{
		Data: map[string]string{"key1": "value1", "key2": "value2"},
	}

	reqBody, err := json.Marshal(reqB)
	assert.NoError(t, err)
	req, err := http.NewRequest("POST", "http://localhost:8081/api/write", bytes.NewBuffer(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	var tuple []byte
	resp.Body.Read(tuple)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, tuple)

	var respData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	assert.NoError(t, err)
	assert.Equal(t, "success", respData["status"])
}

func TestHandleWriteKeys_BadRequest(t *testing.T) {
	client := &http.Client{Timeout: 10 * time.Second}

	reqBody := `{"invalid_json":}`
	req, err := http.NewRequest("POST", "http://localhost:8081/api/write", bytes.NewBufferString(reqBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
