package endpoints_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Andrew-Savin-msk/tarant-kv/internal/config"
	"github.com/Andrew-Savin-msk/tarant-kv/internal/lib/jwt"
)

var token string

type WriteReq struct {
	Data map[string]string `json:"data"`
}

type ReadRequest struct {
	Keys []string `json:"keys"`
}

func TestMain(m *testing.M) {
	cfg := config.Load()

	var err error = nil
	// Генерация JWT токена
	token, err = jwt.GenerateJWT(cfg.UDb.DefUser, time.Hour*24, cfg.Srv.SessionKey)
	if err != nil {
		log.Fatalf("Failed to generate JWT: %v", err)
	}

	reqB := WriteReq{
		Data: map[string]string{"key1": "value1", "key2": "value2"},
	}

	reqBody, err := json.Marshal(reqB)
	if err != nil {
		log.Fatal("Error encoding JSON:", err)
		return
	}
	req, err := http.NewRequest("POST", "http://localhost:8081/api/write", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)
	fmt.Println(err)

	// Выполнение тестов
	code := m.Run()

	// Завершение выполнения тестов с соответствующим кодом
	os.Exit(code)
}
