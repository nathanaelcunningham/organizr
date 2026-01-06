package qbittorrent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Client struct {
	baseURL  string
	username string
	password string
	client   *http.Client
}

func NewClient(baseURL, username, password string) (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create cookie jar: %w", err)
	}
	return &Client{
		baseURL:  strings.TrimSuffix(baseURL, "/"),
		username: username,
		password: password,
		client: &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
		},
	}, nil
}

func (c *Client) Login(ctx context.Context) error {
	data := url.Values{}
	data.Set("username", c.username)
	data.Set("password", c.password)

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v2/auth/login", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("login failed with status %d (failed to read response body: %w)", resp.StatusCode, err)
		}
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
			return fmt.Errorf("authentication failed: invalid username or password")
		}
		return fmt.Errorf("login failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response: %w", err)
	}
	if string(body) != "Ok." {
		return fmt.Errorf("login failed: %s", string(body))
	}

	return nil
}

func (c *Client) AddTorrent(ctx context.Context, magnetLink, torrentURL, category string) (string, error) {
	if err := c.Login(ctx); err != nil {
		return "", fmt.Errorf("failed to authenticate: %w", err)
	}

	data := url.Values{}
	if magnetLink != "" {
		data.Set("urls", magnetLink)
	} else if torrentURL != "" {
		data.Set("urls", torrentURL)
	} else {
		return "", fmt.Errorf("either magnet link or torrent URL must be provided")
	}

	// Add category if provided
	if category != "" {
		data.Set("category", category)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v2/torrents/add", strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create add torrent request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to add torrent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("add torrent failed with status: %d", resp.StatusCode)
	}

	// Extract hash from magnet link or wait for qBittorrent to process
	// For simplicity, extract from magnet link if available
	if magnetLink != "" {
		hash := extractHashFromMagnet(magnetLink)
		if hash != "" {
			return hash, nil
		}
	}

	// If we can't extract hash, we'd need to query the torrent list
	// For now, return empty and let caller handle it
	return "", fmt.Errorf("unable to determine torrent hash")
}

func (c *Client) AddTorrentFromFile(ctx context.Context, torrentData []byte, category string) (string, error) {
	// Validate torrent data
	if len(torrentData) == 0 {
		return "", fmt.Errorf("torrent data is empty")
	}

	// Authenticate first
	if err := c.Login(ctx); err != nil {
		return "", fmt.Errorf("failed to authenticate: %w", err)
	}

	// Create multipart form data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add torrent file
	fileWriter, err := writer.CreateFormFile("torrents", "torrent.torrent")
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := fileWriter.Write(torrentData); err != nil {
		return "", fmt.Errorf("failed to write torrent data: %w", err)
	}

	// Add category if provided
	if category != "" {
		if err := writer.WriteField("category", category); err != nil {
			return "", fmt.Errorf("failed to write category field: %w", err)
		}
	}

	// Close the writer to finalize the multipart message
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create HTTP request with timeout
	uploadCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(uploadCtx, "POST", c.baseURL+"/api/v2/torrents/add", &buf)
	if err != nil {
		return "", fmt.Errorf("failed to create add torrent request: %w", err)
	}

	// Set Content-Type with boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send request
	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to add torrent: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("upload failed with status %d (failed to read error response: %w)", resp.StatusCode, err)
		}
		return "", fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body (should be "Ok." on success, "Fails." on duplicate/error)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read upload response: %w", err)
	}

	responseText := strings.TrimSpace(string(body))

	// Handle different responses
	if responseText != "Ok." && responseText != "Fails." {
		return "", fmt.Errorf("unexpected response from qBittorrent: %s", responseText)
	}

	// "Fails." often means duplicate - we'll check the torrent list to see if it exists
	// "Ok." means successfully added

	// Query torrents to get the hash of the just-added (or existing) torrent
	// Retry up to 3 times as qBittorrent may take 1-2 seconds to process the torrent
	var torrents []TorrentInfo
	maxRetries := 3
	retryDelay := 500 * time.Millisecond

	for attempt := 1; attempt <= maxRetries; attempt++ {
		listReq, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/v2/torrents/info?sort=added_on&reverse=true&limit=10", nil)
		if err != nil {
			return "", fmt.Errorf("failed to create torrent list request: %w", err)
		}

		listResp, err := c.client.Do(listReq)
		if err != nil {
			return "", fmt.Errorf("failed to query torrent list: %w", err)
		}

		if listResp.StatusCode != http.StatusOK {
			listResp.Body.Close()
			return "", fmt.Errorf("torrent list query failed with status: %d", listResp.StatusCode)
		}

		if err := json.NewDecoder(listResp.Body).Decode(&torrents); err != nil {
			listResp.Body.Close()
			return "", fmt.Errorf("failed to decode torrent list: %w", err)
		}
		listResp.Body.Close()

		if len(torrents) > 0 {
			break
		}

		// If not found and not the last attempt, wait before retrying
		if attempt < maxRetries {
			time.Sleep(retryDelay)
		}
	}

	if len(torrents) == 0 {
		if responseText == "Fails." {
			return "", fmt.Errorf("qBittorrent rejected torrent (response: Fails.) - check qBittorrent logs for details. Common causes: invalid save path, disk full, or invalid torrent file")
		}
		return "", fmt.Errorf("torrent not found after upload (tried %d times)", maxRetries)
	}

	// Sort by added time descending (most recent first)
	// Even though we requested sorted in the query, be defensive
	if len(torrents) > 1 {
		sort.Slice(torrents, func(i, j int) bool {
			return torrents[i].AddedOn > torrents[j].AddedOn
		})
	}

	// Return the most recently added torrent's hash
	return torrents[0].Hash, nil
}

func (c *Client) GetTorrentStatus(ctx context.Context, hash string) (string, float64, error) {
	if err := c.Login(ctx); err != nil {
		return "", 0, fmt.Errorf("failed to authenticate: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/v2/torrents/info?hashes="+hash, nil)
	if err != nil {
		return "", 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get torrent info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("get torrent info failed with status: %d", resp.StatusCode)
	}

	var torrents []TorrentInfo
	if err := json.NewDecoder(resp.Body).Decode(&torrents); err != nil {
		return "", 0, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(torrents) == 0 {
		return "", 0, fmt.Errorf("torrent not found")
	}

	torrent := torrents[0]
	status := torrent.State
	progress := torrent.Progress * 100 // Convert to percentage

	return status, progress, nil
}

func (c *Client) GetTorrentFiles(ctx context.Context, hash string) ([]*TorrentFile, error) {
	if err := c.Login(ctx); err != nil {
		return nil, fmt.Errorf("failed to authenticate: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/v2/torrents/files?hash="+hash, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get torrent files: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get torrent files failed with status: %d", resp.StatusCode)
	}

	var filesResp []struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&filesResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Get torrent info to get save path
	infoReq, _ := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/v2/torrents/info?hashes="+hash, nil)
	infoResp, err := c.client.Do(infoReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get torrent info: %w", err)
	}
	defer infoResp.Body.Close()

	var torrents []TorrentInfo
	if err := json.NewDecoder(infoResp.Body).Decode(&torrents); err != nil {
		return nil, fmt.Errorf("failed to decode info response: %w", err)
	}

	if len(torrents) == 0 {
		return nil, fmt.Errorf("torrent not found")
	}

	savePath := torrents[0].SavePath

	files := make([]*TorrentFile, len(filesResp))
	for i, f := range filesResp {
		files[i] = &TorrentFile{
			Name: f.Name,
			Path: savePath + "/" + f.Name,
			Size: f.Size,
		}
	}

	return files, nil
}

func (c *Client) DeleteTorrent(ctx context.Context, hash string, deleteFiles bool) error {
	if err := c.Login(ctx); err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	data := url.Values{}
	data.Set("hashes", hash)
	if deleteFiles {
		data.Set("deleteFiles", "true")
	} else {
		data.Set("deleteFiles", "false")
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/v2/torrents/delete", strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete torrent: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete torrent failed with status: %d", resp.StatusCode)
	}

	return nil
}

func extractHashFromMagnet(magnet string) string {
	// Extract hash from magnet link: magnet:?xt=urn:btih:HASH
	if !strings.HasPrefix(magnet, "magnet:?") {
		return ""
	}

	parts := strings.Split(magnet, "xt=urn:btih:")
	if len(parts) < 2 {
		return ""
	}

	hash := parts[1]
	// Remove anything after the hash (like &dn=name)
	if idx := strings.Index(hash, "&"); idx != -1 {
		hash = hash[:idx]
	}

	return hash
}
