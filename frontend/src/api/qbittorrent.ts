import { api } from './client';

export const qbittorrentApi = {
    testConnection: () =>
        api.get<{ success: boolean; message: string }>('/api/qbittorrent/test'),
};
