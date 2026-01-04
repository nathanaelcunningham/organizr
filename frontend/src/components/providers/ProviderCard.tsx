import React, { useState } from 'react';
import { Card } from '../common/Card';
import { Button } from '../common/Button';
import { Badge } from '../common/Badge';
import type { ProviderConfig } from '../../types/provider';
import { useProviderStore } from '../../stores/useProviderStore';
import { formatDateTime } from '../../utils/formatters';

interface ProviderCardProps {
  provider: ProviderConfig;
  onEdit: () => void;
}

export const ProviderCard: React.FC<ProviderCardProps> = ({
  provider,
  onEdit,
}) => {
  const { deleteProvider, toggleProvider, testConnection } = useProviderStore();
  const [actionLoading, setActionLoading] = useState(false);

  const handleToggle = async () => {
    setActionLoading(true);
    try {
      await toggleProvider(provider.provider_type, !provider.enabled);
    } finally {
      setActionLoading(false);
    }
  };

  const handleTest = async () => {
    setActionLoading(true);
    try {
      await testConnection(provider.provider_type);
    } finally {
      setActionLoading(false);
    }
  };

  const handleDelete = async () => {
    if (
      !window.confirm(
        `Are you sure you want to delete "${provider.display_name}"? This action cannot be undone.`
      )
    ) {
      return;
    }

    setActionLoading(true);
    try {
      await deleteProvider(provider.provider_type);
    } finally {
      setActionLoading(false);
    }
  };

  return (
    <Card>
      <div className="space-y-3">
        {/* Header with name and status */}
        <div className="flex items-start justify-between">
          <div className="flex-1 min-w-0">
            <h3 className="text-lg font-semibold text-gray-900">
              {provider.display_name}
            </h3>
            <p className="text-sm text-gray-500 mt-1">
              Type: {provider.provider_type}
            </p>
          </div>
          <Badge variant={provider.enabled ? 'success' : 'default'} size="md">
            {provider.enabled ? 'Enabled' : 'Disabled'}
          </Badge>
        </div>

        {/* Config Preview (masked) */}
        <div className="bg-gray-50 rounded-lg p-3">
          <p className="text-xs text-gray-500 font-medium mb-2">
            Configuration:
          </p>
          <div className="space-y-1">
            {Object.entries(provider.config).map(([key, value]) => (
              <div key={key} className="flex text-xs">
                <span className="text-gray-600 font-medium w-32 flex-shrink-0">
                  {key}:
                </span>
                <span className="text-gray-800 font-mono truncate">
                  {typeof value === 'string' &&
                  (key.toLowerCase().includes('secret') ||
                    key.toLowerCase().includes('password') ||
                    key.toLowerCase().includes('token'))
                    ? '••••••••'
                    : String(value)}
                </span>
              </div>
            ))}
          </div>
        </div>

        {/* Timestamps */}
        <div className="flex flex-wrap gap-4 text-xs text-gray-500">
          <span>Created: {formatDateTime(provider.created_at)}</span>
          <span>Updated: {formatDateTime(provider.updated_at)}</span>
        </div>

        {/* Action Buttons */}
        <div className="flex flex-wrap gap-2 pt-2">
          <Button
            variant="primary"
            size="sm"
            onClick={handleToggle}
            loading={actionLoading}
          >
            {provider.enabled ? 'Disable' : 'Enable'}
          </Button>
          <Button
            variant="secondary"
            size="sm"
            onClick={handleTest}
            loading={actionLoading}
          >
            Test Connection
          </Button>
          <Button
            variant="secondary"
            size="sm"
            onClick={onEdit}
            disabled={actionLoading}
          >
            Edit
          </Button>
          <Button
            variant="danger"
            size="sm"
            onClick={handleDelete}
            loading={actionLoading}
          >
            Delete
          </Button>
        </div>
      </div>
    </Card>
  );
};
