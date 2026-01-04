export interface AppConfig {
  [key: string]: string;
}

export interface ConfigItem {
  key: string;
  value: string;
}

export interface UpdateConfigRequest {
  value: string;
}

// Common config keys used in the app
export const CONFIG_KEYS = {
  QBITTORRENT_URL: 'qbittorrent.url',
  QBITTORRENT_USERNAME: 'qbittorrent.username',
  QBITTORRENT_PASSWORD: 'qbittorrent.password',
  PATHS_DESTINATION: 'paths.destination',
  PATHS_TEMPLATE: 'paths.template',
  PATHS_NO_SERIES_TEMPLATE: 'paths.no_series_template',
  PATHS_OPERATION: 'paths.operation',
  MONITOR_INTERVAL: 'monitor.interval_seconds',
  MONITOR_AUTO_ORGANIZE: 'monitor.auto_organize',
} as const;
