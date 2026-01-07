package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/nathanael/organizr/internal/models"
)

type MyAnonamouseProvider struct {
	baseUrl string
	secret  string
	client  *http.Client
}

func NewMyAnonamouseProvider(baseUrl, secret string) *MyAnonamouseProvider {
	return &MyAnonamouseProvider{
		baseUrl: baseUrl,
		secret:  secret,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (p *MyAnonamouseProvider) Name() string {
	return "MyAnonamouse"
}

func (p *MyAnonamouseProvider) Search(ctx context.Context, query string) ([]*models.SearchResult, error) {
	// Build search parameters
	params := SearchParams{
		Description: true,
		DLLink:      true,
		ISBN:        false,
		PerPage:     100,
		Torrents: []TorrentSearchParams{
			{
				MainCat:    []MainCategory{CategoryAudiobooks},
				SearchIn:   []SearchIn{SearchInTitle, SearchInAuthor, SearchInSeries, SearchInNarrator},
				SearchType: SearchTypeAll,
				Text:       query,
			},
		},
	}

	// Build URL with parameters
	queryString := formatSearchParamsToURLValues(params)
	searchURL := fmt.Sprintf("%s/tor/js/loadSearchJSONbasic.php?%s", p.baseUrl, queryString)

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("mam_id=%s", p.secret))

	// Make request
	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	// Parse response
	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to SearchResult models
	results := make([]*models.SearchResult, 0, len(searchResp.Data))
	for _, torrent := range searchResp.Data {
		result := &models.SearchResult{
			ID:             fmt.Sprintf("%d", torrent.ID),
			Title:          torrent.Title,
			Author:         formatAuthorInfo(torrent.AuthorInfo),
			Series:         parseSeriesInfo(torrent.SeriesInfo),
			Category:       torrent.CategoryName,
			FileType:       torrent.FileType,
			Language:       torrent.LanguageCode,
			Tags:           parseTags(torrent.Tags),
			Added:          torrent.Added,
			Size:           torrent.Size,
			Seeders:        torrent.Seeders,
			Leechers:       torrent.Leechers,
			NumFiles:       torrent.NumFiles,
			TimesCompleted: torrent.TimesCompleted,
			Freeleech:      torrent.Free == 1,
			FreeleechVIP:   torrent.FLVIP == 1,
			VIP:            torrent.VIP == 1,
			Provider:       p.Name(),
		}

		// Add download link if available
		if torrent.DL != nil && *torrent.DL != "" {
			result.TorrentURL = fmt.Sprintf("%s/tor/download.php?tid=%d", p.baseUrl, torrent.ID)
		}

		// Add description if available
		if torrent.Description != nil {
			result.Description = *torrent.Description
		}

		results = append(results, result)
	}

	return results, nil
}

func (p *MyAnonamouseProvider) DownloadTorrent(ctx context.Context, torrent_id int) ([]byte, error) {
	downloadUrl := fmt.Sprintf("%s/tor/download.php?tid=%d", p.baseUrl, torrent_id)

	req, err := http.NewRequestWithContext(ctx, "GET", downloadUrl, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("mam_id=%s", p.secret))

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {

		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(body))
	}

	respBody, _ := io.ReadAll(resp.Body)

	return respBody, nil
}

func (p *MyAnonamouseProvider) TestConnection(ctx context.Context) error {
	// Make a lightweight search to test authentication
	searchURL := fmt.Sprintf("%s/jsonLoad.php", p.baseUrl)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", fmt.Sprintf("mam_id=%s", p.secret))

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("authentication failed: invalid credentials")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	return nil
}

type SearchParams struct {
	Description bool                  `json:"description"`
	DLLink      bool                  `json:"dl"`
	ISBN        bool                  `json:"isbn"`
	PerPage     int                   `json:"perpage"`
	Torrents    []TorrentSearchParams `json:"tor"`
}

type TorrentSearchParams struct {
	MainCat     []MainCategory `json:"main_cat"`
	SearchIn    []SearchIn     `json:"srchIn"`
	SearchType  SearchType     `json:"searchType"`
	Text        string         `json:"text"`
	StartNumber int            `json:"startNumber"`
	Hash        string         `json:"hash"`
}

type SearchIn string

const (
	SearchInAuthor      SearchIn = "author"
	SearchInDescription SearchIn = "description"
	SearchInFilenames   SearchIn = "filenames"
	SearchInFileTypes   SearchIn = "fileTypes"
	SearchInNarrator    SearchIn = "narrator"
	SearchInSeries      SearchIn = "series"
	SearchInTags        SearchIn = "tags"
	SearchInTitle       SearchIn = "title"
)

type SearchType string

const (
	SearchTypeAll          SearchType = "all"
	SearchTypeActive       SearchType = "active"
	SearchTypeInActive     SearchType = "inactive"
	SearchTypeFreeleech    SearchType = "fl"
	SearchTypeFreeleechVIP SearchType = "fl-VIP"
	SearchTypeVIP          SearchType = "VIP"
	SearchTypeNotVIP       SearchType = "nVIP"
)

type MainCategory int

const (
	CategoryAudiobooks MainCategory = 13
	CategoryBooks      MainCategory = 14
)

type SearchResponse struct {
	Data       []TorrentDetails `json:"data"`
	Total      int              `json:"total"`
	TotalFound int              `json:"total_found"`
}

type TorrentDetails struct {
	ID                int     `json:"id"`
	Added             string  `json:"added"`
	AuthorInfo        string  `json:"author_info"`
	Bookmarked        *string `json:"bookmarked"`
	Category          int     `json:"category"`
	CategoryName      string  `json:"catname"`
	Description       *string `json:"description"`
	DL                *string `json:"dl"`
	FileType          string  `json:"filetype"`
	FLVIP             int     `json:"fl_vip"`
	Free              int     `json:"free"`
	LanguageCode      string  `json:"language_code"`
	Language          int     `json:"language"`
	Leechers          int     `json:"leechers"`
	MainCategory      int     `json:"main_cat"`
	MySnatched        int     `json:"my_snatched"`
	NumFiles          int     `json:"numfiles"`
	OwnerID           int     `json:"owner"`
	OwnerName         string  `json:"owner_name"`
	PersonalFreeleech int     `json:"personal_freeleech"`
	Seeders           int     `json:"seeders"`
	SeriesInfo        string  `json:"series_info"`
	Size              string  `json:"size"`
	Tags              string  `json:"tags"`
	TimesCompleted    int     `json:"times_completed"`
	Title             string  `json:"title"`
	VIP               int     `json:"vip"`
}

func formatAuthorInfo(authorInfo string) string {
	if authorInfo == "" {
		return ""
	}
	authorMap := make(map[string]string)
	if err := json.Unmarshal([]byte(authorInfo), &authorMap); err != nil {
		return ""
	}
	authors := []string{}
	for _, author := range authorMap {
		authors = append(authors, author)
	}

	return strings.Join(authors, ", ")
}

func parseSeriesInfo(seriesInfo string) []models.SeriesInfo {
	if seriesInfo == "" {
		return []models.SeriesInfo{}
	}

	// MAM format: {"123": ["Series Name", "Book Number", numeric_value]}
	seriesMap := make(map[string][]interface{})
	if err := json.Unmarshal([]byte(seriesInfo), &seriesMap); err != nil {
		return []models.SeriesInfo{}
	}

	result := []models.SeriesInfo{}
	for id, s := range seriesMap {
		info := models.SeriesInfo{ID: id}

		if len(s) > 0 {
			if name, ok := s[0].(string); ok {
				info.Name = name
			}
		}
		if len(s) > 1 {
			if number, ok := s[1].(string); ok {
				info.Number = number
			}
		}

		if info.Name != "" {
			result = append(result, info)
		}
	}

	return result
}

func parseTags(tagsStr string) []string {
	if tagsStr == "" {
		return []string{}
	}
	// Tags are comma-separated in the API response
	tags := strings.Split(tagsStr, ",")
	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		if trimmed := strings.TrimSpace(tag); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func formatSearchParamsToURLValues(params SearchParams) string {

	urlParams := ""

	if params.Description {
		urlParams += "description&"
	}
	if params.DLLink {
		urlParams += "dlLink&"
	}
	if params.ISBN {
		urlParams += "isbn&"
	}

	values := url.Values{}

	// Top-level fields
	if params.PerPage > 0 {
		values.Set("perpage", strconv.Itoa(params.PerPage))
	}

	// Nested "tor" array (Torrents)
	for _, torrent := range params.Torrents {
		if len(torrent.MainCat) > 0 {
			for _, category := range torrent.MainCat {
				values.Add("tor[main_cat][]", strconv.Itoa(int(category)))
			}
		}
		if len(torrent.SearchIn) > 0 {
			for _, searchIn := range torrent.SearchIn {
				values.Add(fmt.Sprintf("tor[srchIn][%s]", string(searchIn)), "true")
			}
		}
		if torrent.SearchType != "" {
			values.Set("tor[searchType]", string(torrent.SearchType))
		}
		if torrent.Text != "" {
			values.Set("tor[text]", torrent.Text)
		}
		if torrent.StartNumber > 0 {
			values.Set("tor[startNumber]", strconv.Itoa(torrent.StartNumber))
		}
		if torrent.Hash != "" {
			values.Set("tor[hash]", torrent.Hash)
		}
	}

	return urlParams + values.Encode()
}
