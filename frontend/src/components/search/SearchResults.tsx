import React, { useMemo, useState } from 'react';
import { SeriesGroup } from './SeriesGroup';
import { groupBySeries } from '../../utils/groupSeries';
import { EmptyState } from '../common/EmptyState';
import { Spinner } from '../common/Spinner';
import { Button } from '../common/Button';
import type { SearchResult } from '../../types/search';
import { useDownloadStore } from '../../stores/useDownloadStore';
import type { CreateDownloadRequest } from '../../types/download';

interface SearchResultsProps {
  results: SearchResult[];
  loading: boolean;
  error: string | null;
}

export const SearchResults: React.FC<SearchResultsProps> = ({
  results,
  loading,
  error,
}) => {
  // Group results by series using useMemo for performance
  const grouped = useMemo(() => groupBySeries(results), [results]);

  // Batch selection state
  const [batchMode, setBatchMode] = useState(false);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(new Set());
  const createBatchDownload = useDownloadStore((state) => state.createBatchDownload);
  const [downloadingBatch, setDownloadingBatch] = useState(false);

  // Toggle batch mode
  const toggleBatchMode = () => {
    setBatchMode(!batchMode);
    setSelectedIds(new Set()); // Clear selection when toggling
  };

  // Toggle individual selection
  const toggleSelection = (result: SearchResult) => {
    const id = result.id || result.title; // Use id if available, fallback to title
    setSelectedIds((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(id)) {
        newSet.delete(id);
      } else {
        newSet.add(id);
      }
      return newSet;
    });
  };

  // Check if a result is selected
  const isSelected = (result: SearchResult) => {
    const id = result.id || result.title;
    return selectedIds.has(id);
  };

  // Handle batch download
  const handleBatchDownload = async () => {
    if (selectedIds.size === 0) return;

    setDownloadingBatch(true);
    try {
      // Create download requests from selected results
      const selectedResults = results.filter(r => {
        const id = r.id || r.title;
        return selectedIds.has(id);
      });

      const downloadRequests: CreateDownloadRequest[] = selectedResults.map(result => {
        // Extract first series name and number (books can have multiple series, use primary)
        const series = result.series && result.series.length > 0
          ? result.series[0].name
          : '';
        const seriesNumber = result.series && result.series.length > 0
          ? result.series[0].number
          : '';

        return {
          title: result.title,
          author: result.author,
          series: series,
          seriesNumber: seriesNumber,
          category: 'Audiobooks',
          torrent_url: result.torrent_url,
          magnet_link: result.magnet_link,
        };
      });

      await createBatchDownload(downloadRequests);

      // Clear selection after successful batch
      setSelectedIds(new Set());
    } catch (error) {
      console.error('Failed to create batch downloads:', error);
    } finally {
      setDownloadingBatch(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center py-12">
        <Spinner size="lg" />
        <span className="ml-3 text-gray-600">Searching...</span>
      </div>
    );
  }

  if (error) {
    return (
      <EmptyState
        title="Search Failed"
        description={error}
        icon={
          <svg
            className="w-16 h-16"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        }
      />
    );
  }

  if (results.length === 0) {
    return (
      <EmptyState
        title="No Results Found"
        description="Try adjusting your search query or selecting a different provider"
        icon={
          <svg
            className="w-16 h-16"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        }
      />
    );
  }

  return (
    <div className="relative">
      {/* Header with batch mode toggle */}
      <div className="mb-4 flex items-center justify-between">
        <div className="text-sm text-gray-600">
          Found {results.length} result{results.length === 1 ? '' : 's'} in {grouped.length} {grouped.length === 1 ? 'group' : 'groups'}
        </div>
        <Button
          variant={batchMode ? 'primary' : 'secondary'}
          size="sm"
          onClick={toggleBatchMode}
        >
          {batchMode ? 'Cancel Selection' : 'Select Multiple'}
        </Button>
      </div>

      {/* Results groups */}
      <div>
        {grouped.map((group, index) => (
          <SeriesGroup
            key={`${group.seriesName}-${index}`}
            seriesName={group.seriesName}
            books={group.books}
            batchMode={batchMode}
            selectedIds={selectedIds}
            onToggleSelection={toggleSelection}
            isSelected={isSelected}
          />
        ))}
      </div>

      {/* Floating action bar */}
      {batchMode && selectedIds.size > 0 && (
        <div className="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 shadow-lg transition-transform duration-300 ease-in-out transform translate-y-0 z-50">
          <div className="max-w-7xl mx-auto px-4 py-4 flex items-center justify-between">
            <div className="text-sm font-medium text-gray-700">
              {selectedIds.size} selected
            </div>
            <Button
              variant="primary"
              onClick={handleBatchDownload}
              loading={downloadingBatch}
              disabled={downloadingBatch}
            >
              Download Selected
            </Button>
          </div>
        </div>
      )}
    </div>
  );
};
