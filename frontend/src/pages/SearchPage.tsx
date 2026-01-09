import React, { useCallback } from 'react'
import { PageHeader } from '../components/layout/PageHeader'
import { SearchBar } from '../components/search/SearchBar'
import { SearchResults } from '../components/search/SearchResults'
import { useSearchStore } from '../stores/useSearchStore'

export const SearchPage: React.FC = () => {
  const { loading, error, search, getFilteredResults } = useSearchStore()

  const handleSearch = useCallback(
    (query: string) => {
      search(query)
    },
    [search]
  )

  const filteredResults = getFilteredResults()

  return (
    <div>
      <PageHeader title="Search" subtitle="Search for audiobooks on MyAnonamouse" />
      <SearchBar onSearch={handleSearch} loading={loading} />
      <SearchResults results={filteredResults} loading={loading} error={error} />
    </div>
  )
}
