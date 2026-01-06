import React from 'react';
import { SearchResultListItem } from './SearchResultListItem';
import { EmptyState } from '../common/EmptyState';
import { Spinner } from '../common/Spinner';
import type { SearchResult } from '../../types/search';

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
    <div>
      <div className="mb-4 text-sm text-gray-600">
        Found {results.length} result{results.length === 1 ? '' : 's'}
      </div>
      <div className="divide-y divide-gray-200 border border-gray-200 rounded-lg overflow-hidden">
        {results.map((result, index) => (
          <SearchResultListItem key={`${result.title}-${index}`} result={result} />
        ))}
      </div>
    </div>
  );
};
