import React from 'react';
import type { DownloadStatus } from '../../types/download';

export type BadgeVariant =
  | 'default'
  | 'success'
  | 'error'
  | 'warning'
  | 'info'
  | DownloadStatus;
export type BadgeSize = 'sm' | 'md';

export interface BadgeProps {
  variant?: BadgeVariant;
  size?: BadgeSize;
  children: React.ReactNode;
  className?: string;
}

const variantStyles: Record<string, string> = {
  default: 'bg-gray-100 text-gray-700 border-gray-300',
  success: 'bg-green-100 text-green-700 border-green-300',
  error: 'bg-red-100 text-red-700 border-red-300',
  warning: 'bg-yellow-100 text-yellow-700 border-yellow-300',
  info: 'bg-blue-100 text-blue-700 border-blue-300',
  // Download statuses
  queued: 'bg-gray-100 text-gray-700 border-gray-300',
  downloading: 'bg-blue-100 text-blue-700 border-blue-300',
  completed: 'bg-green-100 text-green-700 border-green-300',
  organizing: 'bg-yellow-100 text-yellow-700 border-yellow-300',
  organized: 'bg-emerald-100 text-emerald-700 border-emerald-300',
  failed: 'bg-red-100 text-red-700 border-red-300',
};

const sizeStyles: Record<BadgeSize, string> = {
  sm: 'px-2 py-0.5 text-xs',
  md: 'px-2.5 py-1 text-sm',
};

export const Badge: React.FC<BadgeProps> = ({
  variant = 'default',
  size = 'sm',
  children,
  className = '',
}) => {
  return (
    <span
      className={`
        inline-flex items-center
        font-medium rounded-full
        border
        ${variantStyles[variant] || variantStyles.default}
        ${sizeStyles[size]}
        ${className}
      `}
    >
      {children}
    </span>
  );
};
