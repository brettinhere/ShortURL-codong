// This Go program implements exactly what main.cod describes.
// Codong compiles to Go — this is what the runtime produces.
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
const dsn = "file:shorturl.db?cache=shared&_journal_mode=WAL"

var db *sql.DB

func generateID(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = chars[rand.Intn(len(chars))]
	}
	return string(b)
}

func isValidURL(s string) bool {
	return strings.HasPrefix(s, "http://") || strings.HasPrefix(s, "https://")
}

func json200(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}

func jsonErr(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var err error
	db, err = sql.Open("sqlite", dsn)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(1)
	db.Exec(`CREATE TABLE IF NOT EXISTS urls (
		code       TEXT PRIMARY KEY,
		long_url   TEXT NOT NULL,
		hits       INTEGER DEFAULT 0,
		created_at TEXT DEFAULT (datetime('now'))
	)`)

	mux := http.NewServeMux()

	// POST /api/shorten
	mux.HandleFunc("POST /api/shorten", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || !isValidURL(body.URL) {
			jsonErr(w, 400, "invalid url")
			return
		}
		code := generateID(6)
		db.Exec("INSERT INTO urls (code, long_url) VALUES (?, ?)", code, body.URL)
		json200(w, map[string]string{
			"code":      code,
			"short_url": "https://codong.org/s/" + code,
		})
	})

	// GET /s/{code}
	mux.HandleFunc("GET /s/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		var longURL string
		if err := db.QueryRow("SELECT long_url FROM urls WHERE code = ?", code).Scan(&longURL); err != nil {
			jsonErr(w, 404, "not found")
			return
		}
		db.Exec("UPDATE urls SET hits = hits + 1 WHERE code = ?", code)
		http.Redirect(w, r, longURL, http.StatusMovedPermanently)
	})

	// GET /api/stats/{code}
	mux.HandleFunc("GET /api/stats/{code}", func(w http.ResponseWriter, r *http.Request) {
		code := r.PathValue("code")
		var row struct {
			Code      string `json:"code"`
			LongURL   string `json:"long_url"`
			Hits      int    `json:"hits"`
			CreatedAt string `json:"created_at"`
		}
		err := db.QueryRow("SELECT code, long_url, hits, created_at FROM urls WHERE code = ?", code).
			Scan(&row.Code, &row.LongURL, &row.Hits, &row.CreatedAt)
		if err != nil {
			jsonErr(w, 404, "not found")
			return
		}
		json200(w, row)
	})

	log.Println("shorturl listening on :8082")
	log.Fatal(http.ListenAndServe(":8082", mux))
}
