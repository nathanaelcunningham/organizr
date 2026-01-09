import React, { useState } from 'react'
import { Input } from '../common/Input'
import { Button } from '../common/Button'
import { MIN_SEARCH_LENGTH } from '../../utils/constants'

interface SearchBarProps {
  onSearch: (query: string) => void
  loading?: boolean
}

export const SearchBar: React.FC<SearchBarProps> = ({ onSearch, loading }) => {
  const [query, setQuery] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (query.length >= MIN_SEARCH_LENGTH) {
      onSearch(query)
    }
  }

  return (
    <form onSubmit={handleSubmit} className="mb-6">
      <div className="flex flex-col sm:flex-row gap-3">
        <div className="flex-1">
          <Input
            type="text"
            placeholder="Search for audiobooks on MyAnonamouse..."
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            className="text-lg"
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
  )
}
