import { api } from './client'
import type { SearchResponse } from '../types/search'

export const searchApi = {
  search: (params: { q: string }) => api.get<SearchResponse>('/api/search', params),

  testConnection: () => api.post<{ success: boolean; message?: string }>('/api/search/test', {}),
}
