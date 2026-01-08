package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/nathanael/organizr/internal/persistence/sqlite"
)

// This tool makes raw API calls to MAM and shows the exact response
// Usage: go run ./cmd/test-mam-raw -query "wheel of time" -db ./organizr.db

func main() {
	query := flag.String("query", "wheel of time", "Search query")
	dbPath := flag.String("db", "./organizr.db", "Path to SQLite database")
	limit := flag.Int("limit", 3, "Number of results to show in detail")
	flag.Parse()

	fmt.Printf("🔍 MAM Raw API Response Inspector\n")
	fmt.Printf("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n\n")

	// Connect to database
	db, err := sqlite.NewDB(*dbPath)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	defer db.Close()

	configRepo := sqlite.NewConfigRepository(db)
	ctx := context.Background()

	// Get MAM credentials
	baseURL, err := configRepo.Get(ctx, "mam.baseurl")
	if err != nil {
		log.Fatalf("❌ MAM base URL not configured: %v", err)
	}
	secret, err := configRepo.Get(ctx, "mam.secret")
	if err != nil {
		log.Fatalf("❌ MAM secret not configured: %v", err)
	}

	fmt.Printf("✓ Configuration loaded\n")
	fmt.Printf("  Base URL: %s\n", baseURL)
	fmt.Printf("  Query: %s\n\n", *query)

	// Build search URL
	searchURL := buildSearchURL(baseURL, *query)
	fmt.Printf("🌐 Request URL:\n%s\n\n", searchURL)

	// Make request
	client := &http.Client{Timeout: 30 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		log.Fatalf("❌ Failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("mam_id=%s", secret))

	fmt.Println("📡 Making API request...")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("❌ Request failed: %v", err)
	}
	defer resp.Body.Close()

	fmt.Printf("✓ Response status: %d\n\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Fatalf("❌ API returned non-200 status: %s", string(body))
	}

	// Read and parse response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("❌ Failed to read response: %v", err)
	}

	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		log.Fatalf("❌ Failed to parse JSON: %v", err)
	}

	fmt.Printf("📊 Response Summary:\n")
	fmt.Printf("   Total results: %d\n", searchResp.Total)
	fmt.Printf("   Total found:   %d\n", len(searchResp.Data))
	fmt.Println()

	// Analyze series_info fields
	seriesCount := 0
	emptySeriesCount := 0
	for _, torrent := range searchResp.Data {
		if torrent.SeriesInfo != "" {
			seriesCount++
		} else {
			emptySeriesCount++
		}
	}

	fmt.Printf("📈 Series Info Statistics:\n")
	fmt.Printf("   With series_info:    %d (%.1f%%)\n", seriesCount, float64(seriesCount)/float64(len(searchResp.Data))*100)
	fmt.Printf("   Without series_info: %d (%.1f%%)\n\n", emptySeriesCount, float64(emptySeriesCount)/float64(len(searchResp.Data))*100)

	// Show detailed results
	fmt.Printf("📚 DETAILED RESULTS (first %d)\n", min(*limit, len(searchResp.Data)))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	for i := 0; i < min(*limit, len(searchResp.Data)); i++ {
		torrent := searchResp.Data[i]

		fmt.Printf("\n[%d] %s\n", i+1, torrent.Title)
		fmt.Printf("─────────────────────────────────────────────\n")
		fmt.Printf("ID:           %d\n", torrent.ID)
		fmt.Printf("Category:     %s\n", torrent.CategoryName)
		fmt.Printf("File Type:    %s\n", torrent.FileType)
		fmt.Printf("Seeders:      %d\n", torrent.Seeders)

		// Show raw author_info
		fmt.Printf("\nRaw author_info:\n")
		if torrent.AuthorInfo != "" {
			var authorData interface{}
			if err := json.Unmarshal([]byte(torrent.AuthorInfo), &authorData); err != nil {
				fmt.Printf("  ⚠️  Invalid JSON: %s\n", torrent.AuthorInfo)
			} else {
				formatted, _ := json.MarshalIndent(authorData, "  ", "  ")
				fmt.Printf("  %s\n", string(formatted))
			}
		} else {
			fmt.Printf("  (empty)\n")
		}

		// Show raw series_info - THIS IS THE KEY FIELD TO DEBUG
		fmt.Printf("\nRaw series_info:\n")
		if torrent.SeriesInfo != "" {
			var seriesData interface{}
			if err := json.Unmarshal([]byte(torrent.SeriesInfo), &seriesData); err != nil {
				fmt.Printf("  ⚠️  Invalid JSON: %s\n", torrent.SeriesInfo)
			} else {
				formatted, _ := json.MarshalIndent(seriesData, "  ", "  ")
				fmt.Printf("  %s\n", string(formatted))

				// Show what our parsing would produce
				parsed := formatSeriesInfo(torrent.SeriesInfo)
				fmt.Printf("\nParsed series (what the app shows):\n")
				fmt.Printf("  \"%s\"\n", parsed)
			}
		} else {
			fmt.Printf("  ⚠️  EMPTY - MAM API returned no series_info for this torrent\n")
		}

		fmt.Println()
	}

	// Show a few examples with series info
	fmt.Println("\n📋 EXAMPLES WITH SERIES INFO")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	count := 0
	for i, torrent := range searchResp.Data {
		if torrent.SeriesInfo != "" && count < 5 {
			fmt.Printf("\n[%d] %s\n", i+1, torrent.Title)
			fmt.Printf("series_info: %s\n", torrent.SeriesInfo)
			fmt.Printf("parsed:      %s\n", formatSeriesInfo(torrent.SeriesInfo))
			count++
		}
	}

	if count == 0 {
		fmt.Println("\n⚠️  NO RESULTS HAD SERIES INFO!")
		fmt.Println("This means either:")
		fmt.Println("1. The search query doesn't match items in a series")
		fmt.Println("2. MAM doesn't have series metadata for these items")
		fmt.Println("3. The API is not returning series_info in the response")
		fmt.Println("\nTry searching for a well-known book series like:")
		fmt.Println("  - \"wheel of time\"")
		fmt.Println("  - \"harry potter\"")
		fmt.Println("  - \"stormlight archive\"")
		fmt.Println("  - \"mistborn\"")
	}

	fmt.Println("\n✅ Inspection complete!")
}

func buildSearchURL(baseURL, query string) string {
	values := url.Values{}
	values.Set("perpage", "100")
	values.Add("tor[main_cat][]", "13") // Audiobooks
	values.Add("tor[srchIn][author]", "true")
	values.Add("tor[srchIn][title]", "true")
	values.Add("tor[srchIn][series]", "true")
	values.Add("tor[srchIn][narrator]", "true")
	values.Set("tor[searchType]", "all")
	values.Set("tor[text]", query)

	return fmt.Sprintf("%s/tor/js/loadSearchJSONbasic.php?description&dlLink&%s", baseURL, values.Encode())
}

func formatSeriesInfo(seriesInfo string) string {
	if seriesInfo == "" {
		return ""
	}
	// MAM returns series info as: {"id": ["Series Name", "Book Number", numeric_value]}
	// The array contains mixed types (strings and numbers), so we use []interface{}
	seriesMap := make(map[string][]interface{})
	if err := json.Unmarshal([]byte(seriesInfo), &seriesMap); err != nil {
		return ""
	}

	series := []string{}
	for _, s := range seriesMap {
		seriesStr := ""
		// First element is the series name (string)
		if len(s) > 0 {
			if name, ok := s[0].(string); ok {
				seriesStr = name
			}
		}
		// Second element is the book number (string)
		if len(s) > 1 && seriesStr != "" {
			if bookNum, ok := s[1].(string); ok && bookNum != "" {
				seriesStr += fmt.Sprintf(" (%s)", bookNum)
			}
		}
		if seriesStr != "" {
			series = append(series, seriesStr)
		}
	}

	result := ""
	for i, s := range series {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Response types
type SearchResponse struct {
	Data       []TorrentDetails `json:"data"`
	Total      int              `json:"total"`
	TotalFound int              `json:"total_found"`
}

type TorrentDetails struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	AuthorInfo   string `json:"author_info"`
	SeriesInfo   string `json:"series_info"`
	CategoryName string `json:"catname"`
	FileType     string `json:"filetype"`
	Seeders      int    `json:"seeders"`
}
