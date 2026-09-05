package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/shivamkrch/olx-clone-api/internal/config"
	"github.com/shivamkrch/olx-clone-api/internal/db"
	"github.com/shivamkrch/olx-clone-api/internal/handlers"
	"github.com/shivamkrch/olx-clone-api/internal/middlewares"
)

func main() {
	cfg := config.MustLoad()

	jsonLogHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	logger := slog.New(jsonLogHandler)
	slog.SetDefault(logger)

	db, err := db.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("main.db.connect: %v", err)
	}

	log.Println("Database connected")

	log.Printf("Starting %s server", cfg.Env)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", handlers.Health)

	lh := handlers.NewListingHandler(db, logger)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("GET /listings/{id}", lh.Get)
	mux.HandleFunc("DELETE /listings/{id}", lh.Delete)

	handler := middlewares.RequestId(mux)
	srv := http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Second * 60,
	}

	log.Printf("%s server is listening on %s", cfg.Env, srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server start failed: %v", err)
	}
}
