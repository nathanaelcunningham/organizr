import { create } from 'zustand';
import { searchApi } from '../api/search';
import type { SearchResult, SearchFilters } from '../types/search';
import { useNotificationStore } from './useNotificationStore';
import { APIClientError } from '../api/client';

interface SearchStore {
  results: SearchResult[];
  loading: boolean;
  error: string | null;
  filters: SearchFilters;

  // Actions
  search: (query: string, provider?: string) => Promise<void>;
  updateFilters: (filters: Partial<SearchFilters>) => void;
  clearResults: () => void;

  // Computed getters
  getFilteredResults: () => SearchResult[];
}

export const useSearchStore = create<SearchStore>((set, get) => ({
  results: [],
  loading: false,
  error: null,
  filters: {
    query: '',
    provider: undefined,
    category: undefined,
    language: undefined,
    minSeeders: undefined,
    freeleechOnly: false,
  },

  search: async (query: string, provider?: string) => {
    if (!query || query.length < 2) {
      set({ results: [], error: 'Query must be at least 2 characters' });
      return;
    }

    try {
      set({ loading: true, error: null });
      const results = await searchApi.search({ q: query, provider });

      // Sort by seeders (descending) by default
      const sortedResults = results.sort((a, b) => b.seeders - a.seeders);

      set({
        results: sortedResults,
        loading: false,
        filters: { ...get().filters, query, provider },
      });
    } catch (error) {
      const message =
        error instanceof APIClientError ? error.message : 'Search failed';
      set({ error: message, loading: false, results: [] });
      useNotificationStore.getState().addNotification('error', message);
    }
  },

  updateFilters: (filters: Partial<SearchFilters>) => {
    set((state) => ({
      filters: { ...state.filters, ...filters },
    }));
  },

  clearResults: () => {
    set({
      results: [],
      error: null,
      filters: {
        query: '',
        provider: undefined,
        category: undefined,
        language: undefined,
        minSeeders: undefined,
        freeleechOnly: false,
      },
    });
  },

  getFilteredResults: () => {
    const { results, filters } = get();

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
  },
}));
