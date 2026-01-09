import type { Download } from '../types/download'
import type { SearchResult } from '../types/search'
import type { AppConfig } from '../types/config'

/**
 * Create a test Download with defaults and optional overrides
 */
export function createTestDownload(overrides?: Partial<Download>): Download {
  const id = overrides?.id || `test-id-${Math.random().toString(36).substring(7)}`

  return {
    id,
    title: 'Test Download',
    author: 'Test Author',
    series: 'Test Series',
    status: 'downloading',
    progress: 50,
    created_at: new Date().toISOString(),
    ...overrides,
  }
}

/**
 * Create a test SearchResult with defaults and optional overrides
 */
export function createTestSearchResult(overrides?: Partial<SearchResult>): SearchResult {
  const id = overrides?.id || `test-result-${Math.random().toString(36).substring(7)}`

  return {
    id,
    title: 'Test Book',
    author: 'Test Author',
    torrent_url: 'https://example.com/torrent.torrent',
    magnet_link: 'magnet:?xt=urn:btih:test',
    provider: 'MyAnonamouse',
    category: 'Audiobooks',
    file_type: 'M4B',
    language: 'English',
    size: '100 MB',
    seeders: 10,
    leechers: 2,
    num_files: 1,
    times_completed: 100,
    freeleech: false,
    freeleech_vip: false,
    vip: false,
    series: [],
    narrator: 'Test Narrator',
    tags: ['fantasy', 'adventure'],
    description: 'Test book description',
    added: '2024-01-01',
    ...overrides,
  }
}

/**
 * Create a test AppConfig with defaults and optional overrides
 */
export function createTestConfig(overrides?: Partial<AppConfig>): AppConfig {
  return {
    'qbittorrent.url': 'http://localhost:8080',
    'qbittorrent.username': 'admin',
    'qbittorrent.password': 'adminpass',
    'mam.baseurl': 'https://www.myanonamouse.net',
    'mam.secret': 'test-secret-key',
    'paths.destination': '/tmp/audiobooks',
    'paths.template': '{author}/{series}/{title}',
    ...overrides,
  }
}
