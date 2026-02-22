package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
)

const Secret = "my_secret_key"

type UrlData struct {
	id  int
	URL string
}

func generateHMAC(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))[:6] // Take the first 8 characters for a shorter key
}

func GetURLData(db *sql.DB, url string) string {
	if len(url) < 7 {
		return fmt.Sprintf("Invalid URL format: %s", url)
	}
	actualId := DecodeBase62(url[6:])
	signature := generateHMAC(url[6:], Secret)
	if !hmac.Equal([]byte(signature), []byte(url[:6])) {
		return fmt.Sprintf("Invalid signature for URL: %s", url)
	}
	var retrievedURL string
	err := db.QueryRow("SELECT long_url FROM urls WHERE id = $1", actualId).Scan(&retrievedURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Sprintf("No data found for URL: %s", url)
		}
		return fmt.Sprintf("Error retrieving data from DB: %v", err)
	}
	return retrievedURL
}

func CreateURLData(db *sql.DB, url string) string {
	var id int
	err := db.QueryRow("INSERT INTO urls (long_url) VALUES ($1) RETURNING id", url).Scan(&id)
	if err != nil {
		return fmt.Sprintf("Error inserting URL into DB: %v", err)
	}
	basedKey := EncodeBase62(id)
	hashedKey := generateHMAC(basedKey, Secret)
	// Simulate data creation based on the URL
	// In a real application, this could involve processing the request body and storing data
	return fmt.Sprintf("Created data for URL: %s link: %s", url, hashedKey+basedKey)
}
