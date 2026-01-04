import { create } from 'zustand';
import { providersApi } from '../api/providers';
import type {
  ProviderType,
  ProviderConfig,
  CreateProviderRequest,
  UpdateProviderRequest,
} from '../types/provider';
import { useNotificationStore } from './useNotificationStore';
import { APIClientError } from '../api/client';

interface ProviderStore {
  providers: ProviderConfig[];
  providerTypes: ProviderType[];
  loading: boolean;
  error: string | null;

  // Actions
  fetchProviders: () => Promise<void>;
  fetchProviderTypes: () => Promise<void>;
  createProvider: (data: CreateProviderRequest) => Promise<boolean>;
  updateProvider: (type: string, data: UpdateProviderRequest) => Promise<boolean>;
  deleteProvider: (type: string) => Promise<void>;
  toggleProvider: (type: string, enabled: boolean) => Promise<void>;
  testConnection: (type: string) => Promise<boolean>;
}

export const useProviderStore = create<ProviderStore>((set) => ({
  providers: [],
  providerTypes: [],
  loading: false,
  error: null,

  fetchProviders: async () => {
    try {
      set({ loading: true, error: null });
      const response = await providersApi.list();
      // Handle both array and object responses
      const providers = Array.isArray(response)
        ? response
        : (response as any)?.providers || [];
      set({ providers, loading: false });
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Failed to fetch providers';
      set({ error: message, loading: false });
      useNotificationStore.getState().addNotification('error', message);
    }
  },

  fetchProviderTypes: async () => {
    try {
      const response = await providersApi.getTypes();
      // Handle both array and object responses
      const providerTypes = Array.isArray(response)
        ? response
        : (response as any)?.provider_types || (response as any)?.providerTypes || [];
      set({ providerTypes });
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Failed to fetch provider types';
      useNotificationStore.getState().addNotification('error', message);
    }
  },

  createProvider: async (data: CreateProviderRequest) => {
    try {
      const provider = await providersApi.create(data);
      set((state) => ({
        providers: [...state.providers, provider],
      }));
      useNotificationStore
        .getState()
        .addNotification('success', 'Provider created successfully');
      return true;
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Failed to create provider';
      useNotificationStore.getState().addNotification('error', message);
      return false;
    }
  },

  updateProvider: async (type: string, data: UpdateProviderRequest) => {
    try {
      const updatedProvider = await providersApi.update(type, data);
      set((state) => ({
        providers: state.providers.map((p) =>
          p.provider_type === type ? updatedProvider : p
        ),
      }));
      useNotificationStore
        .getState()
        .addNotification('success', 'Provider updated successfully');
      return true;
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Failed to update provider';
      useNotificationStore.getState().addNotification('error', message);
      return false;
    }
  },

  deleteProvider: async (type: string) => {
    try {
      await providersApi.delete(type);
      set((state) => ({
        providers: state.providers.filter((p) => p.provider_type !== type),
      }));
      useNotificationStore
        .getState()
        .addNotification('success', 'Provider deleted successfully');
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Failed to delete provider';
      useNotificationStore.getState().addNotification('error', message);
    }
  },

  toggleProvider: async (type: string, enabled: boolean) => {
    try {
      await providersApi.toggle(type, enabled);
      set((state) => ({
        providers: state.providers.map((p) =>
          p.provider_type === type ? { ...p, enabled } : p
        ),
      }));
      useNotificationStore
        .getState()
        .addNotification(
          'success',
          `Provider ${enabled ? 'enabled' : 'disabled'}`
        );
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Failed to toggle provider';
      useNotificationStore.getState().addNotification('error', message);
    }
  },

  testConnection: async (type: string) => {
    try {
      const result = await providersApi.test(type);
      useNotificationStore
        .getState()
        .addNotification(
          result.success ? 'success' : 'error',
          result.message
        );
      return result.success;
    } catch (error) {
      const message =
        error instanceof APIClientError
          ? error.message
          : 'Connection test failed';
      useNotificationStore.getState().addNotification('error', message);
      return false;
    }
  },
}));
