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
  torrent_url?: string;
  magnet_link?: string;
}
