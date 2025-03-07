package repository

import (
	"ShorteningURL/internal/models"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"strings"
	"time"
)

// UrlMapWrapper wrapper for the data model from the table in the DB and *pgxpool.Pool
type UrlMapWrapper struct {
	Url  models.UrlMap
	Pool *pgxpool.Pool
}

// ShortenURL Adds and/or selects url and shortenedUrl
// A non-zero url must be received as input, after which the
// shortened url will be written to the u.Url.ShortenedUrl variable
func (u *UrlMapWrapper) ShortenURL() {

	request := "INSERT INTO url_map (url, shortened_url) VALUES ($1, $2) ON CONFLICT (url) DO NOTHING RETURNING shortened_url;"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u.encodeToShortID()

	if err := u.Pool.QueryRow(ctx, request, u.Url.Url, u.Url.ShortenedUrl).Scan(&u.Url.ShortenedUrl); err != nil {
		log.Printf("Error ShortenURL: %v\n", err)
	}

}

// GetURL returns url by shortenedUrl
// Non-zero shortenedUrl should be received as input,
// after which url will be written to variable u.Url.Url
func (u *UrlMapWrapper) GetURL() {

	request := "SELECT url FROM url_map WHERE shortened_url = $1;"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := u.Pool.QueryRow(ctx, request, u.Url.ShortenedUrl).Scan(&u.Url.Url); err != nil {
		log.Printf("Error ShortenURL: %v\n", err)
	}

}

// DeleteOldURL Deletes all records older than 5 years
func (u *UrlMapWrapper) DeleteOldURL() {

	request := "DELETE FROM url_map WHERE date < NOW() - INTERVAL '5 years'"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := u.Pool.Exec(ctx, request); err != nil {
		log.Printf("Error ShortenURL: %v\n", err)
	}

}

// encodeToShortID - encoding of the link, returns shortenedUrl
func (u *UrlMapWrapper) encodeToShortID() {

	hash := sha256.Sum256([]byte(u.Url.Url))
	hashBase64 := base64.URLEncoding.EncodeToString(hash[:])

	var result strings.Builder
	for _, char := range hashBase64 {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
			result.WriteRune(char)
		}
	}

	u.Url.ShortenedUrl = result.String()
	if len(u.Url.ShortenedUrl) > 5 {
		u.Url.ShortenedUrl = u.Url.ShortenedUrl[:5]
	}

}
