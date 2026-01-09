import React from 'react'
import { Link } from 'react-router-dom'
import { DownloadCard } from './DownloadCard'
import { EmptyState } from '../common/EmptyState'
import { Button } from '../common/Button'
import type { Download, DownloadStatus } from '../../types/download'
import { capitalize } from '../../utils/formatters'

interface DownloadListProps {
  downloads: Download[]
  groupByStatus?: boolean
}

export const DownloadList: React.FC<DownloadListProps> = ({ downloads, groupByStatus = false }) => {
  if (downloads.length === 0) {
    return (
      <EmptyState
        title="No Downloads"
        description="Start by searching for audiobooks and clicking the download button"
        icon={
          <svg className="w-16 h-16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M9 19l3 3m0 0l3-3m-3 3V10"
            />
          </svg>
        }
        action={
          <Link to="/search">
            <Button variant="primary">Search Audiobooks</Button>
          </Link>
        }
      />
    )
  }

  if (!groupByStatus) {
    return (
      <div className="grid gap-4 grid-cols-1 xl:grid-cols-2">
        {downloads.map((download) => (
          <DownloadCard key={download.id} download={download} />
        ))}
      </div>
    )
  }

  // Group downloads by status
  const statusOrder: DownloadStatus[] = [
    'downloading',
    'queued',
    'organizing',
    'completed',
    'organized',
    'failed',
  ]

  const groupedDownloads = statusOrder.reduce(
    (acc, status) => {
      const filtered = downloads.filter((d) => d.status === status)
      if (filtered.length > 0) {
        acc[status] = filtered
      }
      return acc
    },
    {} as Record<DownloadStatus, Download[]>
  )

  return (
    <div className="space-y-6">
      {Object.entries(groupedDownloads).map(([status, statusDownloads]) => (
        <div key={status}>
          <h2 className="text-lg font-semibold text-gray-900 mb-3">
            {capitalize(status)} ({statusDownloads.length})
          </h2>
          <div className="grid gap-4 grid-cols-1 xl:grid-cols-2">
            {statusDownloads.map((download) => (
              <DownloadCard key={download.id} download={download} />
            ))}
          </div>
        </div>
      ))}
    </div>
  )
}
