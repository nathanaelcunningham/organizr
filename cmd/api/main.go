package main

import (
	"context"
	"fmt"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nathanael/organizr/internal/config"
	"github.com/nathanael/organizr/internal/downloads"
	"github.com/nathanael/organizr/internal/persistence/sqlite"
	"github.com/nathanael/organizr/internal/qbittorrent"
	"github.com/nathanael/organizr/internal/search"
	"github.com/nathanael/organizr/internal/server"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// 1. Initialize database
	db, err := sqlite.NewDB("./organizr.db")
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Close()

	// 2. Run migrations
	if err := runMigrations(db); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// 3. Initialize repositories
	downloadRepo := sqlite.NewDownloadRepository(db)
	configRepo := sqlite.NewConfigRepository(db)

	// 4. Initialize config service
	configService := config.NewService(configRepo)

	// 5. Initialize qBittorrent client
	qbURL, err := configRepo.Get(context.Background(), "qbittorrent.url")
	if err != nil {
		qbURL = "http://localhost:8080"
	}
	qbUser, err := configRepo.Get(context.Background(), "qbittorrent.username")
	if err != nil {
		qbUser = "admin"
	}
	qbPass, err := configRepo.Get(context.Background(), "qbittorrent.password")
	if err != nil {
		qbPass = "adminpass"
	}

	qbClient := qbittorrent.NewClient(qbURL, qbUser, qbPass)

	// 6. Initialize services
	downloadService := downloads.NewService(db, qbClient, downloadRepo, configService)
	monitor := downloads.NewMonitor(db, qbClient, downloadRepo, configService)

	// Initialize search providers (users can add their own)
	providers := []search.Provider{
		// User implements their own providers here
		// Example: providers.NewAudiobookBayProvider("https://example.com", "api-key")
	}
	searchService := search.NewService(providers)

	// 7. Start background monitor
	monitorCtx, cancelMonitor := context.WithCancel(context.Background())
	monitorDone := make(chan error, 1)

	go func() {
		if err := monitor.Run(monitorCtx); err != nil && err != context.Canceled {
			monitorDone <- err
		}
	}()

	// 8. Create and start HTTP server
	srv := server.New(server.Config{
		Port:            "8080",
		AllowedOrigins:  []string{"*"},
		DownloadService: downloadService,
		SearchService:   searchService,
		ConfigService:   configService,
	})

	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// 9. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down gracefully...")

	cancelMonitor()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("Shutdown complete")
	return nil
}

func runMigrations(db *sql.DB) error {
	// Read migration file
	migrationSQL, err := os.ReadFile("./assets/migrations/001_init.up.sql")
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	if _, err := db.Exec(string(migrationSQL)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
