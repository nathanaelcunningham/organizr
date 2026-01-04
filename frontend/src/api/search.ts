import { api } from './client';
import type { SearchResult } from '../types/search';

export const searchApi = {
  search: (params: { q: string; provider?: string }) =>
    api.get<SearchResult[]>('/api/search', params),

  getProviders: () =>
    api.get<string[] | { providers: string[] }>('/api/search/providers'),
};
