import { api } from './client'
import type {
  Download,
  CreateDownloadRequest,
  BatchCreateDownloadRequest,
  BatchCreateDownloadResponse,
} from '../types/download'

interface ListDownloadsResponse {
  downloads: Download[]
}

interface GetDownloadResponse {
  download: Download
}

interface CreateDownloadResponse {
  download: Download
}

export const downloadsApi = {
  list: async () => {
    const response = await api.get<ListDownloadsResponse>('/api/downloads')
    return response.downloads
  },

  get: async (id: string) => {
    const response = await api.get<GetDownloadResponse>(`/api/downloads/${id}`)
    return response.download
  },

  create: async (data: CreateDownloadRequest) => {
    const response = await api.post<CreateDownloadResponse>('/api/downloads', data)
    return response.download
  },

  createBatch: async (downloads: CreateDownloadRequest[]) => {
    const request: BatchCreateDownloadRequest = { downloads }
    const response = await api.post<BatchCreateDownloadResponse>('/api/downloads/batch', request)
    return response
  },

  cancel: (id: string) => api.delete<void>(`/api/downloads/${id}`),

  organize: (id: string) => api.post<void>(`/api/downloads/${id}/organize`),
}
