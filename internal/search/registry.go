package search

import (
	"fmt"

	"github.com/nathanael/organizr/internal/models"
	"github.com/nathanael/organizr/internal/search/providers"
)

// ProviderFactory creates a Provider instance from configuration
type ProviderFactory func(config map[string]interface{}) (Provider, error)

// Registry manages provider types and their factories
type Registry struct {
	factories map[string]ProviderFactory
	types     map[string]*models.ProviderType
}

// NewRegistry creates a new registry with built-in providers registered
func NewRegistry() *Registry {
	r := &Registry{
		factories: make(map[string]ProviderFactory),
		types:     make(map[string]*models.ProviderType),
	}

	// Register built-in providers
	r.Register("myanonamouse", myAnonamouseFactory, myAnonamouseType())

	return r
}

// Register adds a provider type to the registry
func (r *Registry) Register(providerType string, factory ProviderFactory, typeInfo *models.ProviderType) {
	r.factories[providerType] = factory
	r.types[providerType] = typeInfo
}

// Create instantiates a provider from configuration
func (r *Registry) Create(providerType string, config map[string]interface{}) (Provider, error) {
	factory, ok := r.factories[providerType]
	if !ok {
		return nil, fmt.Errorf("unknown provider type: %s", providerType)
	}
	return factory(config)
}

// GetTypes returns all available provider types
func (r *Registry) GetTypes() []*models.ProviderType {
	types := make([]*models.ProviderType, 0, len(r.types))
	for _, t := range r.types {
		types = append(types, t)
	}
	return types
}

// GetType returns details about a specific provider type
func (r *Registry) GetType(providerType string) (*models.ProviderType, error) {
	t, ok := r.types[providerType]
	if !ok {
		return nil, fmt.Errorf("unknown provider type: %s", providerType)
	}
	return t, nil
}

// ValidateConfig validates configuration for a provider type
func (r *Registry) ValidateConfig(providerType string, config map[string]interface{}) error {
	typeInfo, err := r.GetType(providerType)
	if err != nil {
		return err
	}

	// Validate required fields
	for _, field := range typeInfo.ConfigSchema {
		if field.Required {
			value, ok := config[field.Name]
			if !ok {
				return fmt.Errorf("missing required field: %s", field.Name)
			}
			// Check if value is empty string
			if strValue, isString := value.(string); isString && strValue == "" {
				return fmt.Errorf("required field cannot be empty: %s", field.Name)
			}
		}
	}

	return nil
}

// Factory functions

func myAnonamouseFactory(config map[string]interface{}) (Provider, error) {
	baseUrl, _ := config["baseUrl"].(string)
	secret, _ := config["secret"].(string)

	if baseUrl == "" {
		baseUrl = "https://www.myanonamouse.net"
	}
	if secret == "" {
		return nil, fmt.Errorf("secret is required")
	}

	return providers.NewMyAnonamouseProvider(baseUrl, secret), nil
}

func myAnonamouseType() *models.ProviderType {
	return &models.ProviderType{
		Type:         "myanonamouse",
		DisplayName:  "MyAnonamouse",
		Description:  "Private audiobook tracker with extensive collection",
		RequiresAuth: true,
		ConfigSchema: []models.ProviderConfigField{
			{
				Name:        "baseUrl",
				DisplayName: "Base URL",
				Type:        "url",
				Required:    false,
				Default:     "https://www.myanonamouse.net",
				Description: "MyAnonamouse base URL",
			},
			{
				Name:        "secret",
				DisplayName: "API Secret",
				Type:        "secret",
				Required:    true,
				Description: "Your MyAnonamouse API secret (mam_id)",
			},
		},
	}
}
