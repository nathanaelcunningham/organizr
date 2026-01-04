import React, { useCallback, useMemo } from 'react';
import { PageHeader } from '../components/layout/PageHeader';
import { SearchBar } from '../components/search/SearchBar';
import { SearchResults } from '../components/search/SearchResults';
import { useSearchStore } from '../stores/useSearchStore';

export const SearchPage: React.FC = () => {
  const { loading, error, search, results, filters } = useSearchStore();

  const handleSearch = useCallback(
    (query: string, provider?: string) => {
      search(query, provider);
    },
    [search]
  );

  const filteredResults = useMemo(() => {
    return results.filter((result) => {
      // Filter by category
      if (filters.category && result.category !== filters.category) {
        return false;
      }

      // Filter by language
      if (filters.language && result.language !== filters.language) {
        return false;
      }

      // Filter by minimum seeders
      if (filters.minSeeders && result.seeders < filters.minSeeders) {
        return false;
      }

      // Filter by freeleech only
      if (filters.freeleechOnly && !result.freeleech) {
        return false;
      }

      return true;
    });
  }, [results, filters]);

  return (
    <div>
      <PageHeader
        title="Search"
        subtitle="Search for audiobooks across configured providers"
      />
      <SearchBar onSearch={handleSearch} loading={loading} />
      <SearchResults results={filteredResults} loading={loading} error={error} />
    </div>
  );
};
