package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"fzsmphone/internal/config"
	"fzsmphone/internal/database"
	"fzsmphone/internal/router"
	"fzsmphone/internal/ws"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (optional, won't error if missing)
	_ = godotenv.Load()

	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[FATAL] Config error: %v", err)
	}

	log.Printf("[INFO] Starting fzsmphone server (env=%s)", cfg.Environment)

	// Connect to database
	ctx := context.Background()
	db, err := database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[FATAL] Database connection error: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.RunMigrations(ctx); err != nil {
		log.Fatalf("[FATAL] Migration error: %v", err)
	}

	// WebSocket hub
	hub := ws.NewHub()
	go hub.Run(ctx)

	// Setup router
	handler := router.New(cfg, db, hub)

	// HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 300 * time.Second, // Must be long enough for AI API streaming (can take 30+s)
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("[INFO] Server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[FATAL] Server error: %v", err)
		}
	}()

	<-done
	log.Println("[INFO] Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("[FATAL] Shutdown error: %v", err)
	}

	log.Println("[INFO] Server stopped")
}
