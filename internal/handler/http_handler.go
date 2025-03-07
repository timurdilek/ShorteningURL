package handler

import (
	"ShorteningURL/internal/models"
	"ShorteningURL/internal/repository"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"log"
	"net/http"
	"time"
)

// Pool stores a variable of type *pgxpool.Pool
// for sending requests to the DB
type Pool struct {
	*pgxpool.Pool
}

// StartServer initializes paths and starts the server
func (p *Pool) StartServer() {

	http.HandleFunc("/get", p.mainPageShortenURL)
	http.HandleFunc("/", p.redirect)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Panic(fmt.Sprintf("Http error: \n%v", err))
	}
}

// mainPageShortenURL The main logic when receiving a POST request, sends a shortened link in the response
func (p *Pool) mainPageShortenURL(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	url := &repository.UrlMapWrapper{Url: models.UrlMap{Url: string(body)}, Pool: p.Pool}
	url.ShortenURL()

	w.Header().Set("Content-Type", "text/plain")

	if _, err = w.Write([]byte(r.Host + "/" + url.Url.ShortenedUrl)); err != nil {
		http.Error(w, "Send error", http.StatusInternalServerError)
	}
}

// redirect Responsible for redirecting users when clicking on a shortened link
func (p *Pool) redirect(w http.ResponseWriter, r *http.Request) {

	url := &repository.UrlMapWrapper{Url: models.UrlMap{ShortenedUrl: r.URL.Path[1:]}, Pool: p.Pool}
	url.GetURL()

	if url.Url.Url == "" {
		http.Error(w, "URL not found", http.StatusNotFound)
	}

	http.Redirect(w, r, url.Url.Url, http.StatusFound)
}

// Ticker removes obsolete records from the DB every 24 hours
// Run in a separate goroutine, since it blocks the execution of the main program
func (p *Pool) Ticker() {

	ticker := time.NewTicker(time.Hour * 24)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			pool := repository.UrlMapWrapper{Pool: p.Pool}
			pool.DeleteOldURL()
			log.Println("Old URLs have been removed.")
		}
	}

}
