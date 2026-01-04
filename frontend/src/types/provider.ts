export interface ProviderConfigField {
  name: string;
  display_name: string;
  type: string; // 'string', 'url', 'secret', 'number', etc.
  required: boolean;
  default?: string;
  description: string;
}

export interface ProviderType {
  type: string; // e.g., "myanonamouse"
  display_name: string; // e.g., "MyAnonamouse"
  description: string;
  requires_auth: boolean;
  config_schema: ProviderConfigField[];
}

export interface ProviderConfig {
  provider_type: string;
  display_name: string;
  enabled: boolean;
  config: Record<string, any>;
  created_at: string;
  updated_at: string;
}

export interface CreateProviderRequest {
  provider_type: string;
  display_name: string;
  enabled: boolean;
  config: Record<string, any>;
}

export interface UpdateProviderRequest {
  display_name: string;
  enabled: boolean;
  config: Record<string, any>;
}

export interface TestConnectionResponse {
  success: boolean;
  message: string;
}
