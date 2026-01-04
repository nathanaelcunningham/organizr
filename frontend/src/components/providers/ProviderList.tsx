import React, { useState } from 'react';
import { Button } from '../common/Button';
import { EmptyState } from '../common/EmptyState';
import { Spinner } from '../common/Spinner';
import { ProviderCard } from './ProviderCard';
import { AddProviderModal } from './AddProviderModal';
import type { ProviderConfig } from '../../types/provider';

interface ProviderListProps {
  providers: ProviderConfig[];
  loading: boolean;
}

export const ProviderList: React.FC<ProviderListProps> = ({
  providers,
  loading,
}) => {
  const [modalOpen, setModalOpen] = useState(false);
  const [editingProvider, setEditingProvider] = useState<
    ProviderConfig | undefined
  >(undefined);

  const handleAdd = () => {
    setEditingProvider(undefined);
    setModalOpen(true);
  };

  const handleEdit = (provider: ProviderConfig) => {
    setEditingProvider(provider);
    setModalOpen(true);
  };

  const handleCloseModal = () => {
    setModalOpen(false);
    setEditingProvider(undefined);
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center py-12">
        <Spinner size="lg" />
        <span className="ml-3 text-gray-600">Loading providers...</span>
      </div>
    );
  }

  if (providers.length === 0) {
    return (
      <>
        <EmptyState
          title="No Providers Configured"
          description="Add a search provider to start finding audiobooks. Providers connect to torrent sites and indexers."
          icon={
            <svg
              className="w-16 h-16"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"
              />
            </svg>
          }
          action={
            <Button variant="primary" onClick={handleAdd}>
              Add Provider
            </Button>
          }
        />
        <AddProviderModal
          isOpen={modalOpen}
          onClose={handleCloseModal}
          existingConfig={editingProvider}
        />
      </>
    );
  }

  return (
    <>
      <div className="mb-4 flex items-center justify-between">
        <p className="text-sm text-gray-600">
          {providers.length} provider{providers.length === 1 ? '' : 's'}{' '}
          configured
        </p>
        <Button variant="primary" onClick={handleAdd}>
          Add Provider
        </Button>
      </div>

      <div className="grid gap-4 grid-cols-1 xl:grid-cols-2">
        {providers.map((provider) => (
          <ProviderCard
            key={provider.provider_type}
            provider={provider}
            onEdit={() => handleEdit(provider)}
          />
        ))}
      </div>

      <AddProviderModal
        isOpen={modalOpen}
        onClose={handleCloseModal}
        existingConfig={editingProvider}
      />
    </>
  );
};
