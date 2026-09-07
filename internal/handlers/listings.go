package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/shivamkrch/olx-clone-api/internal/httpx"
	"github.com/shivamkrch/olx-clone-api/internal/middlewares"
)

type listing struct {
	Id          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListingHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db:     db,
		logger: logger,
	}
}

func (lh *ListingHandler) List(w http.ResponseWriter, r *http.Request) {
	rows, err := lh.db.QueryContext(r.Context(),
		`SELECT id, title, description, price, city, created_at
				FROM listings
				ORDER BY created_at DESC
				LIMIT 100`)
	if err != nil {
		lh.logger.Error("Error querying listings table", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}

	defer rows.Close()

	var listings []listing
	for rows.Next() {
		var l listing

		if err := rows.Scan(&l.Id, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt); err != nil {
			lh.logger.Error("Error scanning listing", "err", err)
			httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
			return
		}

		listings = append(listings, l)
	}

	if err = rows.Err(); err != nil {
		lh.logger.Error("Row iteration error", "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}

	lh.logger.Info("listings fetched", "total", len(listings))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(listings)
}

func (lh *ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.GetRequestIdFromContext(ctx)

	var lc CreateListingRequest
	if err := json.NewDecoder(r.Body).Decode(&lc); err != nil {
		lh.logger.Error("failed to decode request body", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusBadRequest, "Invalid request body.", httpx.CodeMalformedJson)
		return
	}

	query := `INSERT INTO listings (title, description, price, city) VALUES ($1, $2, $3, $4) RETURNING id, title, created_at`
	row := lh.db.QueryRowContext(ctx, query, lc.Title, lc.Description, lc.Price, lc.City)
	var newListing CreateListingResponse
	if err := row.Scan(&newListing.Id, &newListing.Title, &newListing.CreatedAt); err != nil {
		lh.logger.Error("Failed to insert", "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}

	lh.logger.Info("listing created", "request_id", requestId, "listing_id", newListing.Id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(newListing)
}

func (lh *ListingHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")
	listing := lh.getListing(ctx, id)
	if listing == nil {
		httpx.Error(w, http.StatusNotFound, "Requested listing is not available.", httpx.CodeListingNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(listing)
}

func (lh *ListingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middlewares.GetRequestIdFromContext(ctx)

	id := r.PathValue("id")
	if lh.getListing(ctx, id) == nil {
		httpx.Error(w, http.StatusNotFound, "Requested listing is not available.", httpx.CodeListingNotFound)
		return
	}

	_, err := lh.db.ExecContext(ctx, `DELETE FROM listings WHERE id = $1`, id)
	if err != nil {
		lh.logger.Error("delete failed", "listing_id", id, "request_id", requestId, "err", err)
		httpx.Error(w, http.StatusInternalServerError, "Something went wrong", httpx.CodeInternalError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (lh *ListingHandler) getListing(ctx context.Context, id string) *listing {
	requestId := middlewares.GetRequestIdFromContext(ctx)

	row := lh.db.QueryRowContext(ctx, "SELECT id, title, description, price, city, created_at FROM listings WHERE id = $1", id)
	var l listing
	err := row.Scan(&l.Id, &l.Title, &l.Description, &l.Price, &l.City, &l.CreatedAt)
	if err != nil {
		lh.logger.Error("Error scanning listing", "listing_id", id, "request_id", requestId, "err", err)
		return nil
	}

	return &l
}
