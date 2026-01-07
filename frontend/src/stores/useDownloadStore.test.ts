import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useDownloadStore } from './useDownloadStore';
import { downloadsApi } from '../api/downloads';
import type { Download, CreateDownloadRequest } from '../types/download';

// Mock the APIs
vi.mock('../api/downloads', () => ({
  downloadsApi: {
    list: vi.fn(),
    create: vi.fn(),
    cancel: vi.fn(),
    organize: vi.fn(),
  },
}));

// Create a mock function for addNotification that persists across calls
const mockAddNotification = vi.fn();

vi.mock('./useNotificationStore', () => ({
  useNotificationStore: {
    getState: () => ({
      addNotification: mockAddNotification,
    }),
  },
}));

// Mock timers
vi.useFakeTimers();

// Helper to create mock downloads
const createMockDownload = (overrides?: Partial<Download>): Download => ({
  id: '123',
  title: 'Test Book',
  author: 'Test Author',
  series: 'Test Series',
  status: 'downloading',
  progress: 50,
  created_at: new Date().toISOString(),
  ...overrides,
});

describe('useDownloadStore', () => {
  beforeEach(() => {
    // Reset store state before each test
    useDownloadStore.setState({
      downloads: [],
      loading: false,
      error: null,
      pollingInterval: null,
    });
    vi.clearAllMocks();
    mockAddNotification.mockClear();
  });

  afterEach(() => {
    // Clean up any polling intervals
    const { stopPolling } = useDownloadStore.getState();
    stopPolling();
    vi.clearAllTimers();
  });

  describe('fetchDownloads', () => {
    it('should fetch downloads and update state', async () => {
      const mockDownloads = [createMockDownload()];
      vi.mocked(downloadsApi.list).mockResolvedValue(mockDownloads);

      const { fetchDownloads } = useDownloadStore.getState();
      await fetchDownloads();

      const state = useDownloadStore.getState();
      expect(state.downloads).toEqual(mockDownloads);
      expect(state.loading).toBe(false);
      expect(state.error).toBe(null);
    });

    it('should set loading state while fetching', async () => {
      vi.mocked(downloadsApi.list).mockImplementation(
        () =>
          new Promise((resolve) => {
            const state = useDownloadStore.getState();
            expect(state.loading).toBe(true);
            resolve([]);
          })
      );

      const { fetchDownloads } = useDownloadStore.getState();
      await fetchDownloads();
    });

    it('should handle errors and show notification', async () => {
      const error = new Error('Failed to fetch');
      vi.mocked(downloadsApi.list).mockRejectedValue(error);

      const { fetchDownloads } = useDownloadStore.getState();
      await fetchDownloads();

      const state = useDownloadStore.getState();
      expect(state.error).toBeTruthy();
      expect(state.loading).toBe(false);
      expect(mockAddNotification).toHaveBeenCalledWith('error', expect.any(String));
    });
  });

  describe('createDownload', () => {
    it('should create download and add to list', async () => {
      const mockDownload = createMockDownload({ status: 'queued' });
      vi.mocked(downloadsApi.create).mockResolvedValue(mockDownload);

      const request: CreateDownloadRequest = {
        title: 'Test Book',
        author: 'Test Author',
        category: 'Audiobooks',
        magnet_link: 'magnet:...',
      };

      const { createDownload } = useDownloadStore.getState();
      const result = await createDownload(request);

      expect(result).toEqual(mockDownload);
      const state = useDownloadStore.getState();
      expect(state.downloads).toContainEqual(mockDownload);
      expect(mockAddNotification).toHaveBeenCalledWith(
        'success',
        'Download started successfully'
      );
    });

    it('should start polling after creating download', async () => {
      const mockDownload = createMockDownload({ status: 'queued' });
      vi.mocked(downloadsApi.create).mockResolvedValue(mockDownload);

      const request: CreateDownloadRequest = {
        title: 'Test Book',
        author: 'Test Author',
        category: 'Audiobooks',
        magnet_link: 'magnet:...',
      };

      const { createDownload } = useDownloadStore.getState();
      await createDownload(request);

      const state = useDownloadStore.getState();
      expect(state.pollingInterval).not.toBe(null);
    });

    it('should not start polling if already polling', async () => {
      const mockDownload = createMockDownload();
      vi.mocked(downloadsApi.create).mockResolvedValue(mockDownload);

      // Set up existing polling
      useDownloadStore.setState({ pollingInterval: 123 as any });

      const request: CreateDownloadRequest = {
        title: 'Test Book',
        author: 'Test Author',
        category: 'Audiobooks',
        magnet_link: 'magnet:...',
      };

      const { createDownload } = useDownloadStore.getState();
      await createDownload(request);

      const state = useDownloadStore.getState();
      expect(state.pollingInterval).toBe(123);
    });

    it('should handle errors and show notification', async () => {
      const error = new Error('Failed to create');
      vi.mocked(downloadsApi.create).mockRejectedValue(error);

      const request: CreateDownloadRequest = {
        title: 'Test Book',
        author: 'Test Author',
        category: 'Audiobooks',
        magnet_link: 'magnet:...',
      };

      const { createDownload } = useDownloadStore.getState();
      const result = await createDownload(request);

      expect(result).toBeUndefined();
      expect(mockAddNotification).toHaveBeenCalledWith('error', expect.any(String));
    });
  });

  describe('cancelDownload', () => {
    it('should cancel download and remove from list', async () => {
      const mockDownload = createMockDownload();
      useDownloadStore.setState({ downloads: [mockDownload] });
      vi.mocked(downloadsApi.cancel).mockResolvedValue(undefined);

      const { cancelDownload } = useDownloadStore.getState();
      await cancelDownload(mockDownload.id);

      const state = useDownloadStore.getState();
      expect(state.downloads).toHaveLength(0);
      expect(mockAddNotification).toHaveBeenCalledWith('success', 'Download cancelled');
    });

    it('should handle errors and show notification', async () => {
      const error = new Error('Failed to cancel');
      vi.mocked(downloadsApi.cancel).mockRejectedValue(error);

      const { cancelDownload } = useDownloadStore.getState();
      await cancelDownload('123');

      expect(mockAddNotification).toHaveBeenCalledWith('error', expect.any(String));
    });
  });

  describe('organizeDownload', () => {
    it('should organize download and fetch updated status', async () => {
      vi.mocked(downloadsApi.organize).mockResolvedValue(undefined);
      const mockDownloads = [createMockDownload({ status: 'organizing' })];
      vi.mocked(downloadsApi.list).mockResolvedValue(mockDownloads);

      const { organizeDownload } = useDownloadStore.getState();
      await organizeDownload('123');

      expect(downloadsApi.organize).toHaveBeenCalledWith('123');
      expect(mockAddNotification).toHaveBeenCalledWith('success', 'Organizing download...');
      // Should have fetched downloads
      const state = useDownloadStore.getState();
      expect(state.downloads).toEqual(mockDownloads);
    });

    it('should handle errors and show notification', async () => {
      const error = new Error('Failed to organize');
      vi.mocked(downloadsApi.organize).mockRejectedValue(error);

      const { organizeDownload } = useDownloadStore.getState();
      await organizeDownload('123');

      expect(mockAddNotification).toHaveBeenCalledWith('error', expect.any(String));
    });
  });

  describe('polling', () => {
    it('should start polling and fetch downloads at interval', async () => {
      const mockDownloads = [createMockDownload({ status: 'downloading' })];
      vi.mocked(downloadsApi.list).mockResolvedValue(mockDownloads);

      const { startPolling } = useDownloadStore.getState();
      startPolling();

      const state = useDownloadStore.getState();
      expect(state.pollingInterval).not.toBe(null);

      // Advance timer by 3 seconds
      await vi.advanceTimersByTimeAsync(3000);

      expect(downloadsApi.list).toHaveBeenCalled();
    });

    it('should not start polling if already polling', () => {
      const { startPolling } = useDownloadStore.getState();
      startPolling();

      const firstInterval = useDownloadStore.getState().pollingInterval;

      startPolling();

      const secondInterval = useDownloadStore.getState().pollingInterval;
      expect(secondInterval).toBe(firstInterval);
    });

    it('should stop polling when no active downloads', async () => {
      // Start with active downloads
      const activeDownloads = [createMockDownload({ status: 'downloading' })];
      vi.mocked(downloadsApi.list).mockResolvedValueOnce(activeDownloads);

      const { startPolling } = useDownloadStore.getState();
      startPolling();

      await vi.advanceTimersByTimeAsync(3000);
      expect(useDownloadStore.getState().pollingInterval).not.toBe(null);

      // Now return no active downloads
      const completedDownloads = [createMockDownload({ status: 'organized' })];
      vi.mocked(downloadsApi.list).mockResolvedValueOnce(completedDownloads);

      await vi.advanceTimersByTimeAsync(3000);

      // Polling should have stopped
      expect(useDownloadStore.getState().pollingInterval).toBe(null);
    });

    it('should continue polling with queued downloads', async () => {
      const queuedDownloads = [createMockDownload({ status: 'queued' })];
      vi.mocked(downloadsApi.list).mockResolvedValue(queuedDownloads);

      const { startPolling } = useDownloadStore.getState();
      startPolling();

      await vi.advanceTimersByTimeAsync(3000);

      expect(useDownloadStore.getState().pollingInterval).not.toBe(null);
    });

    it('should continue polling with organizing downloads', async () => {
      const organizingDownloads = [createMockDownload({ status: 'organizing' })];
      vi.mocked(downloadsApi.list).mockResolvedValue(organizingDownloads);

      const { startPolling } = useDownloadStore.getState();
      startPolling();

      await vi.advanceTimersByTimeAsync(3000);

      expect(useDownloadStore.getState().pollingInterval).not.toBe(null);
    });

    it('should handle polling errors silently', async () => {
      const error = new Error('Polling failed');
      vi.mocked(downloadsApi.list).mockRejectedValue(error);

      const consoleSpy = vi.spyOn(console, 'error').mockImplementation(() => {});

      const { startPolling } = useDownloadStore.getState();
      startPolling();

      await vi.advanceTimersByTimeAsync(3000);

      expect(consoleSpy).toHaveBeenCalledWith('Polling error:', error);
      // Should not call notification store during polling errors
      expect(mockAddNotification).not.toHaveBeenCalled();

      consoleSpy.mockRestore();
    });

    it('should stop polling on stopPolling call', () => {
      const { startPolling, stopPolling } = useDownloadStore.getState();
      startPolling();

      expect(useDownloadStore.getState().pollingInterval).not.toBe(null);

      stopPolling();

      expect(useDownloadStore.getState().pollingInterval).toBe(null);
    });

    it('should not error when stopPolling called without active polling', () => {
      const { stopPolling } = useDownloadStore.getState();
      expect(() => stopPolling()).not.toThrow();
    });
  });

  describe('computed getters', () => {
    beforeEach(() => {
      const downloads = [
        createMockDownload({ id: '1', status: 'queued' }),
        createMockDownload({ id: '2', status: 'downloading' }),
        createMockDownload({ id: '3', status: 'organizing' }),
        createMockDownload({ id: '4', status: 'completed' }),
        createMockDownload({ id: '5', status: 'organized' }),
        createMockDownload({ id: '6', status: 'failed' }),
      ];
      useDownloadStore.setState({ downloads });
    });

    it('getActiveDownloads should return queued, downloading, and organizing', () => {
      const { getActiveDownloads } = useDownloadStore.getState();
      const active = getActiveDownloads();

      expect(active).toHaveLength(3);
      expect(active.map((d) => d.status)).toEqual([
        'queued',
        'downloading',
        'organizing',
      ]);
    });

    it('getCompletedDownloads should return only completed', () => {
      const { getCompletedDownloads } = useDownloadStore.getState();
      const completed = getCompletedDownloads();

      expect(completed).toHaveLength(1);
      expect(completed[0].status).toBe('completed');
    });

    it('getFailedDownloads should return only failed', () => {
      const { getFailedDownloads } = useDownloadStore.getState();
      const failed = getFailedDownloads();

      expect(failed).toHaveLength(1);
      expect(failed[0].status).toBe('failed');
    });
  });
});
