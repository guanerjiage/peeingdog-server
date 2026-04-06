package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"

	"peeingdog-server/config"
	"peeingdog-server/db"
	"peeingdog-server/handlers"
	"peeingdog-server/sql/queries/generated"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Verify database connection
	if err := database.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("✓ Connected to PostgreSQL database")

	// Initialize sqlc queries
	queries := generated.New(database)

	// Initialize handlers with sqlc queries
	h := handlers.New(queries)

	// Setup routes
	mux := http.NewServeMux()
	
	// Health check
	mux.HandleFunc("GET /health", h.Health)
	
	// User endpoints
	mux.HandleFunc("GET /api/users", h.GetUsers)
	mux.HandleFunc("POST /api/users", h.CreateUser)
	mux.HandleFunc("GET /api/users/{id}", h.GetUser)
	
	// Message endpoints
	mux.HandleFunc("POST /api/users/{id}/messages", h.CreateMessage)
	mux.HandleFunc("GET /api/users/{id}/messages", h.GetUserMessages)
	mux.HandleFunc("GET /api/messages/nearby", h.GetNearbyMessages)

	// Create HTTP server
	addr := fmt.Sprintf(":%d", cfg.Port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Use errgroup to manage concurrent goroutines
	eg, ctx := errgroup.WithContext(context.Background())

	// Channel for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine 1: Run the server
	eg.Go(func() error {
		log.Printf("Server starting on http://localhost:%d\n", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
		return nil
	})

	// Goroutine 2: Listen for shutdown signals or context cancellation
	eg.Go(func() error {
		select {
		case sig := <-sigChan:
			log.Printf("\nReceived signal: %v", sig)
		case <-ctx.Done():
			log.Println("Context cancelled, shutting down...")
			return ctx.Err()
		}

		// Graceful shutdown with 10 second timeout
		log.Println("Shutting down server gracefully...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown error: %w", err)
		}

		return nil
	})

	// Wait for all goroutines to complete
	if err := eg.Wait(); err != nil {
		log.Printf("Error: %v", err)
	}

	log.Println("✓ Server stopped")
}
