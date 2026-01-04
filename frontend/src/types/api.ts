export interface APIError {
  error: string;
  message: string;
  code?: number;
}

export interface APIResponse<T = any> {
  data?: T;
  error?: APIError;
}

export interface HealthStatus {
  status: string;
  database: string;
  qbittorrent: string;
  monitor: string;
}
