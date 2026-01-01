package qbittorrent

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL  string
	username string
	password string
	client   *http.Client
}

func NewClient(baseURL, username, password string) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		baseURL:  strings.TrimSuffix(baseURL, "/"),
		username: username,
		password: password,
		client: &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
		},
	}
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
		return fmt.Errorf("login failed with status: %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "Ok." {
		return fmt.Errorf("login failed: %s", string(body))
	}

	return nil
}

func (c *Client) AddTorrent(ctx context.Context, magnetLink, torrentURL string) (string, error) {
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
