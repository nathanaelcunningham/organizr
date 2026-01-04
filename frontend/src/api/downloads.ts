import { api } from './client';
import type { Download, CreateDownloadRequest } from '../types/download';

export const downloadsApi = {
  list: () => api.get<Download[]>('/api/downloads'),

  get: (id: string) => api.get<Download>(`/api/downloads/${id}`),

  create: (data: CreateDownloadRequest) =>
    api.post<Download>('/api/downloads', data),

  cancel: (id: string) => api.delete<void>(`/api/downloads/${id}`),

  organize: (id: string) => api.post<void>(`/api/downloads/${id}/organize`),
};
