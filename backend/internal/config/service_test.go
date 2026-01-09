package config

import (
	"context"
	"os"
	"testing"

	"github.com/nathanael/organizr/internal/testutil"
)

// mockConfigRepository is a mock implementation of persistence.ConfigRepository
type mockConfigRepository struct {
	getFunc    func(ctx context.Context, key string) (string, error)
	getAllFunc func(ctx context.Context) (map[string]string, error)
	setFunc    func(ctx context.Context, key, value string) error
}

func (m *mockConfigRepository) Get(ctx context.Context, key string) (string, error) {
	if m.getFunc != nil {
		return m.getFunc(ctx, key)
	}
	return "", nil
}

func (m *mockConfigRepository) GetAll(ctx context.Context) (map[string]string, error) {
	if m.getAllFunc != nil {
		return m.getAllFunc(ctx)
	}
	return make(map[string]string), nil
}

func (m *mockConfigRepository) Set(ctx context.Context, key, value string) error {
	if m.setFunc != nil {
		return m.setFunc(ctx, key, value)
	}
	return nil
}

func TestService_Get_EnvironmentPrecedence(t *testing.T) {
	tests := []struct {
		name          string
		key           string
		envKey        string
		envValue      string
		dbValue       string
		expectedValue string
	}{
		{
			name:          "environment variable overrides database",
			key:           "qbittorrent.url",
			envKey:        "QBITTORRENT_URL",
			envValue:      "http://env-qbit:8080",
			dbValue:       "http://db-qbit:8080",
			expectedValue: "http://env-qbit:8080",
		},
		{
			name:          "falls back to database when env not set",
			key:           "qbittorrent.url",
			envKey:        "QBITTORRENT_URL",
			envValue:      "",
			dbValue:       "http://db-qbit:8080",
			expectedValue: "http://db-qbit:8080",
		},
		{
			name:          "qbittorrent username from env",
			key:           "qbittorrent.username",
			envKey:        "QBITTORRENT_USERNAME",
			envValue:      "envuser",
			dbValue:       "dbuser",
			expectedValue: "envuser",
		},
		{
			name:          "paths.destination from env",
			key:           "paths.destination",
			envKey:        "PATHS_DESTINATION",
			envValue:      "/env/audiobooks",
			dbValue:       "/db/audiobooks",
			expectedValue: "/env/audiobooks",
		},
		{
			name:          "monitor.auto_organize from env",
			key:           "monitor.auto_organize",
			envKey:        "MONITOR_AUTO_ORGANIZE",
			envValue:      "false",
			dbValue:       "true",
			expectedValue: "false",
		},
		{
			name:          "unknown config key uses database only",
			key:           "unknown.key",
			envKey:        "",
			envValue:      "",
			dbValue:       "db-value",
			expectedValue: "db-value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			if tt.envValue != "" && tt.envKey != "" {
				os.Setenv(tt.envKey, tt.envValue)
				defer os.Unsetenv(tt.envKey)
			}

			// Mock repository
			mockRepo := &mockConfigRepository{
				getFunc: func(ctx context.Context, key string) (string, error) {
					return tt.dbValue, nil
				},
			}

			svc := NewService(mockRepo)
			value, err := svc.Get(context.Background(), tt.key)

			testutil.AssertNoError(t, err)
			testutil.AssertEqual(t, value, tt.expectedValue)
		})
	}
}

func TestService_GetAll_MergesEnvironmentVariables(t *testing.T) {
	tests := []struct {
		name            string
		envVars         map[string]string
		dbConfigs       map[string]string
		expectedConfigs map[string]string
	}{
		{
			name: "environment overrides database values",
			envVars: map[string]string{
				"QBITTORRENT_URL":  "http://env-qbit:8080",
				"PATHS_DESTINATION": "/env/audiobooks",
			},
			dbConfigs: map[string]string{
				"qbittorrent.url":   "http://db-qbit:8080",
				"qbittorrent.username": "admin",
				"paths.destination": "/db/audiobooks",
			},
			expectedConfigs: map[string]string{
				"qbittorrent.url":      "http://env-qbit:8080",
				"qbittorrent.username": "admin",
				"paths.destination":    "/env/audiobooks",
			},
		},
		{
			name: "database values used when env not set",
			envVars: map[string]string{},
			dbConfigs: map[string]string{
				"qbittorrent.url":      "http://db-qbit:8080",
				"qbittorrent.username": "admin",
				"paths.destination":    "/db/audiobooks",
			},
			expectedConfigs: map[string]string{
				"qbittorrent.url":      "http://db-qbit:8080",
				"qbittorrent.username": "admin",
				"paths.destination":    "/db/audiobooks",
			},
		},
		{
			name: "all config options from environment",
			envVars: map[string]string{
				"QBITTORRENT_URL":          "http://qbit:8080",
				"QBITTORRENT_USERNAME":     "user",
				"QBITTORRENT_PASSWORD":     "pass",
				"PATHS_DESTINATION":        "/audiobooks",
				"PATHS_TEMPLATE":           "{author}/{series}/{title}",
				"PATHS_NO_SERIES_TEMPLATE": "{author}/{title}",
				"PATHS_OPERATION":          "move",
				"PATHS_LOCAL_MOUNT":        "/mnt/data",
				"MONITOR_INTERVAL_SECONDS": "60",
				"MONITOR_AUTO_ORGANIZE":    "true",
				"MAM_BASEURL":              "https://mam.net",
				"MAM_SECRET":               "secret123",
			},
			dbConfigs: map[string]string{
				"qbittorrent.url": "http://default:8080",
			},
			expectedConfigs: map[string]string{
				"qbittorrent.url":          "http://qbit:8080",
				"qbittorrent.username":     "user",
				"qbittorrent.password":     "pass",
				"paths.destination":        "/audiobooks",
				"paths.template":           "{author}/{series}/{title}",
				"paths.no_series_template": "{author}/{title}",
				"paths.operation":          "move",
				"paths.local_mount":        "/mnt/data",
				"monitor.interval_seconds": "60",
				"monitor.auto_organize":    "true",
				"mam.baseurl":              "https://mam.net",
				"mam.secret":               "secret123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			for key, value := range tt.envVars {
				os.Setenv(key, value)
				defer os.Unsetenv(key)
			}

			// Mock repository
			mockRepo := &mockConfigRepository{
				getAllFunc: func(ctx context.Context) (map[string]string, error) {
					return tt.dbConfigs, nil
				},
			}

			svc := NewService(mockRepo)
			configs, err := svc.GetAll(context.Background())

			testutil.AssertNoError(t, err)
			testutil.AssertEqual(t, configs, tt.expectedConfigs)
		})
	}
}

func TestService_Set_WritesToDatabaseOnly(t *testing.T) {
	// Set environment variable
	os.Setenv("QBITTORRENT_URL", "http://env-qbit:8080")
	defer os.Unsetenv("QBITTORRENT_URL")

	setCalled := false
	mockRepo := &mockConfigRepository{
		setFunc: func(ctx context.Context, key, value string) error {
			setCalled = true
			testutil.AssertEqual(t, key, "qbittorrent.url")
			testutil.AssertEqual(t, value, "http://new-qbit:8080")
			return nil
		},
		getFunc: func(ctx context.Context, key string) (string, error) {
			// Return database value
			return "http://db-qbit:8080", nil
		},
	}

	svc := NewService(mockRepo)

	// Set should write to database only, not modify environment
	err := svc.Set(context.Background(), "qbittorrent.url", "http://new-qbit:8080")
	testutil.AssertNoError(t, err)

	if !setCalled {
		t.Error("expected Set to call repository")
	}

	// Get should still return environment value (not the new database value)
	value, err := svc.Get(context.Background(), "qbittorrent.url")
	testutil.AssertNoError(t, err)
	testutil.AssertEqual(t, value, "http://env-qbit:8080")
}
