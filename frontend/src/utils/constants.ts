import type { DownloadStatus } from '../types/download';

/**
 * Status color mappings for badges and progress bars
 */
export const STATUS_COLORS: Record<DownloadStatus, string> = {
  queued: 'bg-gray-100 text-gray-700 border-gray-300',
  downloading: 'bg-blue-100 text-blue-700 border-blue-300',
  completed: 'bg-green-100 text-green-700 border-green-300',
  organizing: 'bg-yellow-100 text-yellow-700 border-yellow-300',
  organized: 'bg-emerald-100 text-emerald-700 border-emerald-300',
  failed: 'bg-red-100 text-red-700 border-red-300',
};

/**
 * Status labels for display
 */
export const STATUS_LABELS: Record<DownloadStatus, string> = {
  queued: 'Queued',
  downloading: 'Downloading',
  completed: 'Completed',
  organizing: 'Organizing',
  organized: 'Organized',
  failed: 'Failed',
};

/**
 * Active download statuses (used for polling)
 */
export const ACTIVE_STATUSES: DownloadStatus[] = [
  'queued',
  'downloading',
  'organizing',
];

/**
 * Polling interval for downloads (milliseconds)
 */
export const DOWNLOAD_POLL_INTERVAL = 3000;

/**
 * Search debounce delay (milliseconds)
 */
export const SEARCH_DEBOUNCE_DELAY = 500;

/**
 * Minimum search query length
 */
export const MIN_SEARCH_LENGTH = 2;
