import { api } from './client';
import type { AppConfigResponse, UpdateConfigRequest } from '../types/config';

export interface PreviewPathRequest {
    template: string;
    author: string;
    series?: string;
    title: string;
}

export interface PreviewPathResponse {
    valid: boolean;
    path?: string;
    error?: string;
}

export const configApi = {
    getAll: () => api.get<AppConfigResponse>('/api/config'),

    get: (key: string) => api.get<{ value: string }>(`/api/config/${key}`),

    update: (key: string, data: UpdateConfigRequest) =>
        api.put<void>(`/api/config/${key}`, data),

    previewPath: (data: PreviewPathRequest) =>
        api.post<PreviewPathResponse>('/api/config/preview-path', data),
};
