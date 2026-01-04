import { create } from 'zustand';
import { downloadsApi } from '../api/downloads';
import type { Download, CreateDownloadRequest } from '../types/download';
import { useNotificationStore } from './useNotificationStore';
import { APIClientError } from '../api/client';

interface DownloadStore {
  downloads: Download[];
  loading: boolean;
  error: string | null;
  pollingInterval: number | null;

  // Actions
  fetchDownloads: () => Promise<void>;
  createDownload: (data: CreateDownloadRequest) => Promise<Download | undefined>;
  cancelDownload: (id: string) => Promise<void>;
  organizeDownload: (id: string) => Promise<void>;
  startPolling: () => void;
  stopPolling: () => void;

  // Computed getters
  getActiveDownloads: () => Download[];
  getCompletedDownloads: () => Download[];
  getFailedDownloads: () => Download[];
}

export const useDownloadStore = create<DownloadStore>((set, get) => ({
  downloads: [],
  loading: false,
  error: null,
  pollingInterval: null,

  fetchDownloads: async () => {
    try {
      set({ loading: true, error: null });
      const downloads = await downloadsApi.list();
      set({ downloads, loading: false });
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Failed to fetch downloads';
      set({ error: message, loading: false });
      useNotificationStore.getState().addNotification('error', message);
    }
  },

  createDownload: async (data: CreateDownloadRequest) => {
    try {
      const download = await downloadsApi.create(data);
      set((state) => ({
        downloads: [download, ...state.downloads],
      }));
      useNotificationStore
        .getState()
        .addNotification('success', 'Download started successfully');

      // Start polling if not already polling
      const { pollingInterval, startPolling } = get();
      if (!pollingInterval) {
        startPolling();
      }

      return download;
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Failed to create download';
      useNotificationStore.getState().addNotification('error', message);
      return undefined;
    }
  },

  cancelDownload: async (id: string) => {
    try {
      await downloadsApi.cancel(id);
      set((state) => ({
        downloads: state.downloads.filter((d) => d.id !== id),
      }));
      useNotificationStore
        .getState()
        .addNotification('success', 'Download cancelled');
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Failed to cancel download';
      useNotificationStore.getState().addNotification('error', message);
    }
  },

  organizeDownload: async (id: string) => {
    try {
      await downloadsApi.organize(id);
      useNotificationStore
        .getState()
        .addNotification('success', 'Organizing download...');
      // Fetch updated status
      await get().fetchDownloads();
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Failed to organize download';
      useNotificationStore.getState().addNotification('error', message);
    }
  },

  startPolling: () => {
    const { pollingInterval, stopPolling } = get();

    // Don't start if already polling
    if (pollingInterval) return;

    const intervalId = window.setInterval(async () => {
      // Fetch downloads without showing loading state
      try {
        const downloads = await downloadsApi.list();
        set({ downloads });

        // Check if there are any active downloads
        const hasActive = downloads.some((d) =>
          ['queued', 'downloading', 'organizing'].includes(d.status)
        );

        // Stop polling if no active downloads
        if (!hasActive) {
          stopPolling();
        }
      } catch (error) {
        // Silently fail during polling to avoid spamming notifications
        console.error('Polling error:', error);
      }
    }, 3000); // Poll every 3 seconds

    set({ pollingInterval: intervalId });
  },

  stopPolling: () => {
    const { pollingInterval } = get();
    if (pollingInterval) {
      clearInterval(pollingInterval);
      set({ pollingInterval: null });
    }
  },

  // Computed getters
  getActiveDownloads: () =>
    get().downloads.filter((d) =>
      ['queued', 'downloading', 'organizing'].includes(d.status)
    ),

  getCompletedDownloads: () =>
    get().downloads.filter((d) => d.status === 'completed'),

  getFailedDownloads: () =>
    get().downloads.filter((d) => d.status === 'failed'),
}));
