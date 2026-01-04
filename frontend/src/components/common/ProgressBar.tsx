import React from 'react';
import type { DownloadStatus } from '../../types/download';

export type ProgressBarSize = 'sm' | 'md' | 'lg';

export interface ProgressBarProps {
  progress: number; // 0-100
  status?: DownloadStatus;
  showLabel?: boolean;
  size?: ProgressBarSize;
  className?: string;
}

const statusColors: Record<DownloadStatus, string> = {
  queued: 'bg-gray-500',
  downloading: 'bg-blue-600',
  completed: 'bg-green-600',
  organizing: 'bg-yellow-600',
  organized: 'bg-emerald-600',
  failed: 'bg-red-600',
};

const sizeStyles: Record<ProgressBarSize, string> = {
  sm: 'h-1',
  md: 'h-2',
  lg: 'h-3',
};

export const ProgressBar: React.FC<ProgressBarProps> = ({
  progress,
  status = 'downloading',
  showLabel = false,
  size = 'md',
  className = '',
}) => {
  const normalizedProgress = Math.max(0, Math.min(100, progress));

  return (
    <div className={`w-full ${className}`}>
      {showLabel && (
        <div className="flex justify-between items-center mb-1">
          <span className="text-sm font-medium text-gray-700 capitalize">
            {status}
          </span>
          <span className="text-sm font-medium text-gray-700">
            {normalizedProgress}%
          </span>
        </div>
      )}
      <div className={`w-full bg-gray-200 rounded-full overflow-hidden ${sizeStyles[size]}`}>
        <div
          className={`${sizeStyles[size]} rounded-full transition-all duration-300 ${statusColors[status]}`}
          style={{ width: `${normalizedProgress}%` }}
        />
      </div>
    </div>
  );
};
