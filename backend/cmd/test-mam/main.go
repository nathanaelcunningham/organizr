package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nathanael/organizr/internal/persistence/sqlite"
	"github.com/nathanael/organizr/internal/search"
)

// Simple CLI tool to test MAM search with real API calls
// Usage: go run ./cmd/test-mam -query "wheel of time" -db ./organizr.db

func main() {
	// Parse flags
	query := flag.String("query", "wheel of time", "Search query")
	dbPath := flag.String("db", "./organizr.db", "Path to SQLite database")
	debug := flag.Bool("debug", false, "Enable debug output")
	flag.Parse()

	fmt.Printf("🔍 Testing MAM Search API\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Connect to database
	fmt.Printf("📁 Connecting to database: %s\n", *dbPath)
	db, err := sqlite.NewDB(*dbPath)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()
	fmt.Println("✓ Database connected")

	// Create repositories
	configRepo := sqlite.NewConfigRepository(db)

	// Verify MAM configuration
	ctx := context.Background()
	baseURL, err := configRepo.Get(ctx, "mam.baseurl")
	if err != nil {
		log.Fatalf("❌ MAM base URL not configured: %v", err)
	}
	secret, err := configRepo.Get(ctx, "mam.secret")
	if err != nil {
		log.Fatalf("❌ MAM secret not configured: %v", err)
	}

	fmt.Printf("✓ MAM configuration loaded\n")
	fmt.Printf("  Base URL: %s\n", baseURL)
	fmt.Printf("  Secret: %s...\n\n", secret[:10])

	// Create MAM service
	mamService := search.NewMAMService(configRepo)

	// Test connection
	fmt.Println("🔌 Testing connection to MAM...")
	if err := mamService.TestConnection(ctx); err != nil {
		log.Fatalf("❌ Failed to connect to MAM: %v", err)
	}
	fmt.Println("✓ Successfully connected to MAM")

	// Perform search
	fmt.Printf("🔍 Searching for: \"%s\"\n\n", *query)
	results, err := mamService.Search(ctx, *query)
	if err != nil {
		log.Fatalf("❌ Search failed: %v", err)
	}

	fmt.Printf("Found %d results\n", len(results))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Analyze results
	seriesCount := 0
	emptySeriesCount := 0
	resultsWithSeries := []int{}

	for i, result := range results {
		if len(result.Series) > 0 {
			seriesCount++
			resultsWithSeries = append(resultsWithSeries, i)
		} else {
			emptySeriesCount++
		}
	}

	// Print summary
	fmt.Println("📊 SUMMARY")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Total results:          %d\n", len(results))
	fmt.Printf("With series info:       %d (%.1f%%)\n", seriesCount, float64(seriesCount)/float64(len(results))*100)
	fmt.Printf("Without series info:    %d (%.1f%%)\n\n", emptySeriesCount, float64(emptySeriesCount)/float64(len(results))*100)

	// Print first 5 results in detail
	fmt.Println("📚 FIRST 5 RESULTS")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	for i := 0; i < min(5, len(results)); i++ {
		result := results[i]
		fmt.Printf("\n#%d: %s\n", i+1, result.Title)
		fmt.Printf("   ID:       %s\n", result.ID)
		fmt.Printf("   Author:   %s\n", result.Author)

		if len(result.Series) > 0 {
			// Format series info for display
			seriesStrs := make([]string, len(result.Series))
			for i, s := range result.Series {
				if s.Number != "" {
					seriesStrs[i] = fmt.Sprintf("%s (%s)", s.Name, s.Number)
				} else {
					seriesStrs[i] = s.Name
				}
			}
			fmt.Printf("   Series:   ✓ %s\n", strings.Join(seriesStrs, ", "))
		} else {
			fmt.Printf("   Series:   ⚠️  EMPTY\n")
		}

		fmt.Printf("   Category: %s\n", result.Category)
		fmt.Printf("   FileType: %s\n", result.FileType)
		fmt.Printf("   Seeders:  %d\n", result.Seeders)
		if result.Freeleech {
			fmt.Printf("   Special:  🎁 FREELEECH\n")
		}
	}

	// If debug mode, show results with series
	if *debug && len(resultsWithSeries) > 0 {
		fmt.Println("\n\n🔍 DEBUG: RESULTS WITH SERIES INFO")
		fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		for _, idx := range resultsWithSeries {
			if idx >= 10 { // Only show first 10
				break
			}
			result := results[idx]
			fmt.Printf("\n#%d: %s\n", idx+1, result.Title)
			fmt.Printf("   Author: %s\n", result.Author)
			fmt.Printf("   Series: %s\n", result.Series)
		}
	}

	// Save results to file if requested
	if os.Getenv("SAVE_RESULTS") == "1" {
		filename := "mam_search_results.json"
		data, err := json.MarshalIndent(results, "", "  ")
		if err != nil {
			log.Printf("⚠️  Failed to marshal results: %v", err)
		} else {
			if err := os.WriteFile(filename, data, 0644); err != nil {
				log.Printf("⚠️  Failed to write results to file: %v", err)
			} else {
				fmt.Printf("\n💾 Results saved to: %s\n", filename)
			}
		}
	}

	fmt.Println("\n✅ Test complete!")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
