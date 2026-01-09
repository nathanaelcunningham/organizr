import React from 'react'
import type { DownloadStatus } from '../../types/download'

export type FilterStatus = 'all' | DownloadStatus

interface DownloadFiltersProps {
  activeFilter: FilterStatus
  onFilterChange: (filter: FilterStatus) => void
  counts: Record<FilterStatus, number>
}

const filters: { value: FilterStatus; label: string }[] = [
  { value: 'all', label: 'All' },
  { value: 'queued', label: 'Queued' },
  { value: 'downloading', label: 'Downloading' },
  { value: 'completed', label: 'Completed' },
  { value: 'organizing', label: 'Organizing' },
  { value: 'organized', label: 'Organized' },
  { value: 'failed', label: 'Failed' },
]

export const DownloadFilters: React.FC<DownloadFiltersProps> = ({
  activeFilter,
  onFilterChange,
  counts,
}) => {
  return (
    <div className="mb-6 border-b border-gray-200 -mx-4 sm:mx-0 px-4 sm:px-0">
      <nav className="-mb-px flex gap-4 sm:gap-6 overflow-x-auto scrollbar-hide">
        {filters.map((filter) => {
          const count = counts[filter.value] || 0
          const isActive = activeFilter === filter.value

          return (
            <button
              key={filter.value}
              onClick={() => onFilterChange(filter.value)}
              className={`
                py-3 px-1 border-b-2 font-medium text-sm whitespace-nowrap
                transition-colors
                ${
                  isActive
                    ? 'border-blue-600 text-blue-600'
                    : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                }
              `}
            >
              {filter.label}
              {count > 0 && (
                <span
                  className={`
                  ml-2 py-0.5 px-2 rounded-full text-xs font-semibold
                  ${isActive ? 'bg-blue-100 text-blue-600' : 'bg-gray-100 text-gray-600'}
                `}
                >
                  {count}
                </span>
              )}
            </button>
          )
        })}
      </nav>
    </div>
  )
}
