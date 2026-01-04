export interface SearchResult {
  id?: string;
  title: string;
  author: string;
  series?: string;
  narrator?: string;
  torrent_url?: string;
  magnet_link?: string;
  provider: string;
  size: string;
  seeders: number;
  leechers: number;
  category?: string;
  file_type?: string;
  language?: string;
  tags?: string[];
  description?: string;
  freeleech?: boolean;
  freeleech_vip?: boolean;
  vip?: boolean;
  num_files?: number;
  times_completed?: number;
  added?: string;
}

export interface SearchFilters {
  query: string;
  provider?: string;
  category?: string;
  language?: string;
  minSeeders?: number;
  freeleechOnly?: boolean;
}
