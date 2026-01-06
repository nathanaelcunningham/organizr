package main

import (
	"context"
	"database/sql"
	"fmt"
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

	// 5. Initialize MAM service
	mamService := search.NewMAMService(configRepo)

	// 6. Initialize qBittorrent client
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

	// 7. Initialize download services
	downloadService := downloads.NewService(db, qbClient, downloadRepo, configService, mamService)
	monitor := downloads.NewMonitor(db, qbClient, downloadRepo, configService)

	// 8. Start background monitor
	monitorCtx, cancelMonitor := context.WithCancel(context.Background())
	monitorDone := make(chan error, 1)

	go func() {
		if err := monitor.Run(monitorCtx); err != nil && err != context.Canceled {
			monitorDone <- err
		}
	}()

	// 9. Create and start HTTP server
	srv := server.New(server.Config{
		Port:            "8080",
		AllowedOrigins:  []string{"*"},
		DownloadService: downloadService,
		SearchService:   mamService,
		ConfigService:   configService,
	})

	go func() {
		if err := srv.Start(); err != nil {
			log.Printf("Server error: %v", err)
		}
	}()

	// 10. Graceful shutdown
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
	// Create migrations tracking table
	createMigrationsTable := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`
	if _, err := db.Exec(createMigrationsTable); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Define migrations in order
	migrations := []struct {
		version  int
		filename string
	}{
		{1, "./assets/migrations/001_init.up.sql"},
		{2, "./assets/migrations/002_add_category.up.sql"},
	}

	for _, migration := range migrations {
		// Check if already applied
		var count int
		err := db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", migration.version).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration status: %w", err)
		}

		if count > 0 {
			log.Printf("Migration %d already applied, skipping", migration.version)
			continue
		}

		// Read and execute migration
		migrationSQL, err := os.ReadFile(migration.filename)
		if err != nil {
			return fmt.Errorf("failed to read migration %d: %w", migration.version, err)
		}

		if _, err := db.Exec(string(migrationSQL)); err != nil {
			return fmt.Errorf("failed to execute migration %d: %w", migration.version, err)
		}

		// Record migration
		if _, err := db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", migration.version); err != nil {
			return fmt.Errorf("failed to record migration %d: %w", migration.version, err)
		}

		log.Printf("Applied migration %d", migration.version)
	}

	log.Println("All migrations completed successfully")
	return nil
}
