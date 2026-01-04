import React, { useEffect, useState } from 'react';
import { PageHeader } from '../components/layout/PageHeader';
import { DownloadFilters } from '../components/downloads/DownloadFilters';
import type { FilterStatus } from '../components/downloads/DownloadFilters';
import { DownloadList } from '../components/downloads/DownloadList';
import { useDownloadStore } from '../stores/useDownloadStore';

export const DownloadsPage: React.FC = () => {
  const { downloads, fetchDownloads, startPolling, stopPolling } =
    useDownloadStore();
  const [activeFilter, setActiveFilter] = useState<FilterStatus>('all');

  useEffect(() => {
    fetchDownloads();
    startPolling();

    return () => {
      stopPolling();
    };
  }, [fetchDownloads, startPolling, stopPolling]);

  // Filter downloads based on active filter
  const filteredDownloads =
    activeFilter === 'all'
      ? downloads
      : downloads.filter((d) => d.status === activeFilter);

  // Calculate counts for each filter
  const counts: Record<FilterStatus, number> = {
    all: downloads.length,
    queued: downloads.filter((d) => d.status === 'queued').length,
    downloading: downloads.filter((d) => d.status === 'downloading').length,
    completed: downloads.filter((d) => d.status === 'completed').length,
    organizing: downloads.filter((d) => d.status === 'organizing').length,
    organized: downloads.filter((d) => d.status === 'organized').length,
    failed: downloads.filter((d) => d.status === 'failed').length,
  };

  return (
    <div>
      <PageHeader
        title="Downloads"
        subtitle="Manage your audiobook downloads"
      />
      <DownloadFilters
        activeFilter={activeFilter}
        onFilterChange={setActiveFilter}
        counts={counts}
      />
      <DownloadList
        downloads={filteredDownloads}
        groupByStatus={activeFilter === 'all'}
      />
    </div>
  );
};
