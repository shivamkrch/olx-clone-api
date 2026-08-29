package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type listing struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"createdAt"`
}

func List(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		rows, err := db.QueryContext(ctx,
			`SELECT id, title, description, price, city, created_at
				FROM listings
				ORDER BY created_at DESC
				LIMIT 100`)
		if err != nil {
			log.Printf("Error querying listings table: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		defer rows.Close()

		var listings []listing
		for rows.Next() {
			var l listing

			if err := rows.Scan(&l.Id, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
				log.Printf("Error scanning listing: %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			listings = append(listings, l)
		}

		if err = rows.Err(); err != nil {
			log.Printf("Row iteration error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(listings)
	}
}
