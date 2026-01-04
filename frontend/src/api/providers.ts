import { api } from './client';
import type {
  ProviderType,
  ProviderConfig,
  CreateProviderRequest,
  UpdateProviderRequest,
  TestConnectionResponse,
} from '../types/provider';

export const providersApi = {
  getTypes: () =>
    api.get<
      ProviderType[] | { provider_types: ProviderType[] } | { providerTypes: ProviderType[] }
    >('/api/search/providers/types'),

  list: () =>
    api.get<ProviderConfig[] | { providers: ProviderConfig[] }>(
      '/api/search/providers/config'
    ),

  get: (type: string) =>
    api.get<ProviderConfig>(`/api/search/providers/config/${type}`),

  create: (data: CreateProviderRequest) =>
    api.post<ProviderConfig>('/api/search/providers/config', data),

  update: (type: string, data: UpdateProviderRequest) =>
    api.put<ProviderConfig>(`/api/search/providers/config/${type}`, data),

  delete: (type: string) =>
    api.delete<void>(`/api/search/providers/config/${type}`),

  toggle: (type: string, enabled: boolean) =>
    api.patch<void>(`/api/search/providers/config/${type}/toggle`, {
      enabled,
    }),

  test: (type: string) =>
    api.post<TestConnectionResponse>(
      `/api/search/providers/config/${type}/test`
    ),
};
