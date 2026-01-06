import { api } from './client';
import type { AppConfigResponse, UpdateConfigRequest } from '../types/config';

export const configApi = {
    getAll: () => api.get<AppConfigResponse>('/api/config'),

    get: (key: string) => api.get<{ value: string }>(`/api/config/${key}`),

    update: (key: string, data: UpdateConfigRequest) =>
        api.put<void>(`/api/config/${key}`, data),
};
