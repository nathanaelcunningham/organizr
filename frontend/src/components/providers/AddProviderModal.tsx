import React, { useState } from 'react';
import { Modal } from '../common/Modal';
import { Select } from '../common/Select';
import { ProviderForm } from './ProviderForm';
import { useProviderStore } from '../../stores/useProviderStore';
import type { ProviderConfig } from '../../types/provider';

interface AddProviderModalProps {
  isOpen: boolean;
  onClose: () => void;
  existingConfig?: ProviderConfig;
}

export const AddProviderModal: React.FC<AddProviderModalProps> = ({
  isOpen,
  onClose,
  existingConfig,
}) => {
  const { providerTypes, createProvider, updateProvider, testConnection } =
    useProviderStore();

  const [selectedType, setSelectedType] = useState<string>(
    existingConfig?.provider_type || ''
  );

  const selectedProviderType = providerTypes.find(
    (type) => type.type === selectedType
  );

  const handleSubmit = async (data: {
    display_name: string;
    enabled: boolean;
    config: Record<string, string>;
  }) => {
    if (existingConfig) {
      // Update existing provider
      const success = await updateProvider(existingConfig.provider_type, {
        display_name: data.display_name,
        enabled: data.enabled,
        config: data.config,
      });
      if (success) {
        onClose();
      }
      return success;
    } else {
      // Create new provider
      const success = await createProvider({
        provider_type: selectedType,
        display_name: data.display_name,
        enabled: data.enabled,
        config: data.config,
      });
      if (success) {
        onClose();
      }
      return success;
    }
  };

  const handleTest = async () => {
    if (existingConfig) {
      return await testConnection(existingConfig.provider_type);
    }
    // For new providers, we can't test until they're created
    // Could implement a test endpoint that accepts config without saving
    return false;
  };

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={existingConfig ? 'Edit Provider' : 'Add Provider'}
      size="lg"
    >
      <div className="space-y-4">
        {!existingConfig && (
          <Select
            label="Provider Type"
            value={selectedType}
            onChange={(e) => setSelectedType(e.target.value)}
            options={[
              { value: '', label: 'Select a provider type...' },
              ...providerTypes.map((type) => ({
                value: type.type,
                label: type.display_name,
              })),
            ]}
            required
          />
        )}

        {selectedProviderType && (
          <ProviderForm
            providerType={selectedProviderType}
            existingConfig={existingConfig}
            onSubmit={handleSubmit}
            onCancel={onClose}
            onTest={existingConfig ? handleTest : undefined}
          />
        )}

        {!selectedProviderType && !existingConfig && (
          <p className="text-sm text-gray-500 text-center py-8">
            Select a provider type to continue
          </p>
        )}
      </div>
    </Modal>
  );
};
