export type DownloadStatus =
    | 'queued'
    | 'downloading'
    | 'completed'
    | 'organizing'
    | 'organized'
    | 'failed';

export interface Download {
    id: string;
    title: string;
    author: string;
    series?: string;
    seriesNumber?: string;
    status: DownloadStatus;
    progress: number; // 0-100
    organized_path?: string;
    error_message?: string;
    created_at: string;
    completed_at?: string;
    organized_at?: string;
}

export interface CreateDownloadRequest {
    title: string;
    author: string;
    series?: string;
    seriesNumber?: string;
    category: string;
    torrent_url?: string;
    magnet_link?: string;
}

export interface BatchCreateDownloadRequest {
    downloads: CreateDownloadRequest[];
}

export interface BatchDownloadError {
    index: number;
    request: CreateDownloadRequest;
    error: string;
}

export interface BatchCreateDownloadResponse {
    successful: Download[];
    failed: BatchDownloadError[];
}
