import React, { useState } from 'react';
import { Card } from '../common/Card';
import { Button } from '../common/Button';
import { Badge } from '../common/Badge';
import { ProgressBar } from '../common/ProgressBar';
import type { Download } from '../../types/download';
import { useDownloadStore } from '../../stores/useDownloadStore';
import { useNotificationStore } from '../../stores/useNotificationStore';
import { formatRelativeTime, capitalize } from '../../utils/formatters';

interface DownloadCardProps {
  download: Download;
}

export const DownloadCard: React.FC<DownloadCardProps> = ({ download }) => {
  const { cancelDownload, organizeDownload } = useDownloadStore();
  const { addNotification } = useNotificationStore();
  const [actionLoading, setActionLoading] = useState(false);
  const [copied, setCopied] = useState(false);

  const handleCancel = async () => {
    if (
      !window.confirm(
        'Are you sure you want to cancel this download? This action cannot be undone.'
      )
    ) {
      return;
    }
    setActionLoading(true);
    try {
      await cancelDownload(download.id);
    } finally {
      setActionLoading(false);
    }
  };

  const handleOrganize = async () => {
    setActionLoading(true);
    try {
      await organizeDownload(download.id);
    } finally {
      setActionLoading(false);
    }
  };

  const handleCopyPath = async () => {
    if (!download.organized_path) return;

    try {
      await navigator.clipboard.writeText(download.organized_path);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      addNotification('error', 'Failed to copy path to clipboard');
    }
  };

  const showCancelButton =
    download.status === 'queued' ||
    download.status === 'downloading' ||
    download.status === 'organizing';

  const showOrganizeButton = download.status === 'completed';
  const showRetryButton =
    download.status === 'failed' && download.error_message?.includes('organiz');

  return (
    <Card className="hover:shadow-md transition-shadow">
      <div className="space-y-3">
        {/* Title and Author */}
        <div className="flex items-start justify-between">
          <div className="flex-1 min-w-0">
            <h3 className="text-lg font-semibold text-gray-900 truncate">
              {download.title}
            </h3>
            <p className="text-sm text-gray-600 mt-1">by {download.author}</p>
            {download.series && (
              <p className="text-sm text-gray-500">Series: {download.series}</p>
            )}
          </div>
          <Badge
            variant={download.status}
            size="md"
            className={download.status === 'organizing' ? 'animate-pulse' : ''}
          >
            {capitalize(download.status)}
          </Badge>
        </div>

        {/* Progress Bar for Active Downloads */}
        {download.status === 'downloading' && (
          <ProgressBar
            progress={download.progress}
            status={download.status}
            showLabel
            size="md"
          />
        )}

        {/* Organization in Progress */}
        {download.status === 'organizing' && (
          <div className="space-y-2">
            <ProgressBar progress={100} status="organizing" showLabel={false} size="md" />
            <p className="text-sm text-blue-700">
              Organizing files... Creating folder structure and copying files to destination.
            </p>
          </div>
        )}

        {/* Organized Path */}
        {download.status === 'organized' && download.organized_path && (
          <div className="bg-emerald-50 border border-emerald-200 rounded-lg p-3">
            <div className="flex items-center justify-between mb-1">
              <p className="text-xs text-emerald-700 font-medium">
                Organized Path:
              </p>
              <button
                onClick={handleCopyPath}
                className="text-xs text-emerald-700 hover:text-emerald-800 font-medium px-2 py-1 rounded hover:bg-emerald-100 transition-colors"
                title="Copy path to clipboard"
              >
                {copied ? '✓ Copied!' : '📋 Copy'}
              </button>
            </div>
            <p className="text-sm text-emerald-900 font-mono break-all">
              {download.organized_path}
            </p>
          </div>
        )}

        {/* Error Message */}
        {download.status === 'failed' && download.error_message && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-3">
            <p className="text-xs text-red-700 font-medium mb-1">Error:</p>
            <p className="text-sm text-red-900">{download.error_message}</p>
          </div>
        )}

        {/* Timestamps */}
        <div className="flex flex-wrap gap-4 text-xs text-gray-500">
          <span>Created {formatRelativeTime(download.created_at)}</span>
          {download.completed_at && (
            <span>Completed {formatRelativeTime(download.completed_at)}</span>
          )}
          {download.organized_at && (
            <span>Organized {formatRelativeTime(download.organized_at)}</span>
          )}
        </div>

        {/* Action Buttons */}
        {(showCancelButton || showOrganizeButton || showRetryButton) && (
          <div className="flex gap-2 pt-2">
            {showOrganizeButton && (
              <Button
                variant="primary"
                size="sm"
                onClick={handleOrganize}
                loading={actionLoading}
                className="flex-1"
              >
                Organize Now
              </Button>
            )}
            {showRetryButton && (
              <Button
                variant="primary"
                size="sm"
                onClick={handleOrganize}
                loading={actionLoading}
                className="flex-1"
              >
                Retry Organization
              </Button>
            )}
            {showCancelButton && (
              <Button
                variant="danger"
                size="sm"
                onClick={handleCancel}
                loading={actionLoading}
                className={showOrganizeButton || showRetryButton ? 'w-auto' : 'flex-1'}
              >
                Cancel
              </Button>
            )}
          </div>
        )}
      </div>
    </Card>
  );
};
