import { create } from 'zustand';
import { configApi } from '../api/config';
import type { AppConfig } from '../types/config';
import { useNotificationStore } from './useNotificationStore';
import { APIClientError } from '../api/client';

interface ConfigStore {
  config: AppConfig;
  loading: boolean;
  error: string | null;

  // Actions
  fetchConfig: () => Promise<void>;
  updateConfig: (key: string, value: string) => Promise<boolean>;
  updateMultipleConfigs: (updates: Record<string, string>) => Promise<boolean>;

  // Helper getters
  getConfigValue: (key: string, defaultValue?: string) => string | undefined;
}

export const useConfigStore = create<ConfigStore>((set, get) => ({
  config: {},
  loading: false,
  error: null,

  fetchConfig: async () => {
    try {
      set({ loading: true, error: null });
      const config = await configApi.getAll();
      set({ config, loading: false });
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Failed to fetch configuration';
      set({ error: message, loading: false });
      useNotificationStore.getState().addNotification('error', message);
    }
  },

  updateConfig: async (key: string, value: string) => {
    try {
      await configApi.update(key, { value });
      set((state) => ({
        config: { ...state.config, [key]: value },
      }));
      useNotificationStore
        .getState()
        .addNotification('success', 'Configuration updated successfully');
      return true;
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Failed to update configuration';
      useNotificationStore.getState().addNotification('error', message);
      return false;
    }
  },

  updateMultipleConfigs: async (updates: Record<string, string>) => {
    try {
      // Update all configs in parallel
      await Promise.all(
        Object.entries(updates).map(([key, value]) =>
          configApi.update(key, { value })
        )
      );

      set((state) => ({
        config: { ...state.config, ...updates },
      }));

      useNotificationStore
        .getState()
        .addNotification('success', 'Configuration updated successfully');
      return true;
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Failed to update configuration';
      useNotificationStore.getState().addNotification('error', message);
      return false;
    }
  },

  getConfigValue: (key: string, defaultValue?: string) => {
    const { config } = get();
    return config[key] ?? defaultValue;
  },
}));
