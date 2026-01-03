package models

import "time"

// ProviderConfig represents a configured provider instance in the database
type ProviderConfig struct {
	ProviderType string                 // Primary key: "myanonamouse", "audiobookbay"
	DisplayName  string                 // Human-readable name
	Enabled      bool                   // Enable/disable toggle
	ConfigJSON   map[string]interface{} // Provider-specific config
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// ProviderType describes an available provider type and its configuration requirements
type ProviderType struct {
	Type         string                // Unique identifier
	DisplayName  string                // Human-readable name
	Description  string                // Description for UI
	ConfigSchema []ProviderConfigField // Required/optional fields
	RequiresAuth bool                  // Whether authentication is needed
}

// ProviderConfigField defines a configuration field for a provider type
type ProviderConfigField struct {
	Name        string // Field name (e.g., "baseUrl", "secret")
	DisplayName string // Human-readable label
	Type        string // "string", "url", "secret", etc.
	Required    bool   // Whether field is required
	Default     string // Default value if any
	Description string // Help text
}
