import React, { useState, useEffect } from 'react';
import { Input } from '../common/Input';
import { Select } from '../common/Select';
import { Button } from '../common/Button';
import { useDebounce } from '../../hooks/useDebounce';
import { searchApi } from '../../api/search';
import { MIN_SEARCH_LENGTH, SEARCH_DEBOUNCE_DELAY } from '../../utils/constants';

interface SearchBarProps {
  onSearch: (query: string, provider?: string) => void;
  loading?: boolean;
}

export const SearchBar: React.FC<SearchBarProps> = ({ onSearch, loading }) => {
  const [query, setQuery] = useState('');
  const [provider, setProvider] = useState<string>('all');
  const [providers, setProviders] = useState<string[]>([]);
  const debouncedQuery = useDebounce(query, SEARCH_DEBOUNCE_DELAY);

  // Fetch available providers on mount
  useEffect(() => {
    const fetchProviders = async () => {
      try {
        const response = await searchApi.getProviders();
        // Handle both array and object responses
        const providerList = Array.isArray(response)
          ? response
          : (response as any)?.providers || [];
        setProviders(providerList);
      } catch (error) {
        console.error('Failed to fetch providers:', error);
        setProviders([]); // Ensure providers is always an array
      }
    };
    fetchProviders();
  }, []);

  // Auto-search when debounced query changes
  useEffect(() => {
    if (debouncedQuery.length >= MIN_SEARCH_LENGTH) {
      onSearch(debouncedQuery, provider === 'all' ? undefined : provider);
    }
  }, [debouncedQuery, provider, onSearch]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (query.length >= MIN_SEARCH_LENGTH) {
      onSearch(query, provider === 'all' ? undefined : provider);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="mb-6">
      <div className="flex flex-col sm:flex-row gap-3">
        <div className="flex-1">
          <Input
            type="text"
            placeholder="Search for audiobooks..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="text-lg"
          />
        </div>
        <div className="w-full sm:w-48">
          <Select
            value={provider}
            onChange={(e) => setProvider(e.target.value)}
            options={[
              { value: 'all', label: 'All Providers' },
              ...providers.map((p) => ({ value: p, label: p })),
            ]}
          />
        </div>
        <Button
          type="submit"
          variant="primary"
          size="md"
          disabled={query.length < MIN_SEARCH_LENGTH || loading}
          loading={loading}
          className="w-full sm:w-auto"
        >
          Search
        </Button>
      </div>
      {query.length > 0 && query.length < MIN_SEARCH_LENGTH && (
        <p className="mt-2 text-sm text-gray-500">
          Enter at least {MIN_SEARCH_LENGTH} characters to search
        </p>
      )}
    </form>
  );
};
