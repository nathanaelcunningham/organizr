import React, { useState } from 'react';
import { Button } from '../common/Button';
import { Input } from '../common/Input';
import { ProviderFormField } from './ProviderFormField';
import type { ProviderType, ProviderConfig } from '../../types/provider';

interface ProviderFormProps {
  providerType: ProviderType;
  existingConfig?: ProviderConfig;
  onSubmit: (data: {
    display_name: string;
    enabled: boolean;
    config: Record<string, string>;
  }) => Promise<boolean>;
  onCancel: () => void;
  onTest?: (config: Record<string, string>) => Promise<boolean>;
}

export const ProviderForm: React.FC<ProviderFormProps> = ({
  providerType,
  existingConfig,
  onSubmit,
  onCancel,
  onTest,
}) => {
  const [displayName, setDisplayName] = useState(
    existingConfig?.display_name || providerType.display_name
  );
  const [enabled, setEnabled] = useState(existingConfig?.enabled ?? true);
  const [config, setConfig] = useState<Record<string, string>>(() => {
    const initial: Record<string, string> = {};
    providerType.config_schema.forEach((field) => {
      initial[field.name] =
        existingConfig?.config[field.name] || field.default || '';
    });
    return initial;
  });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitting, setSubmitting] = useState(false);
  const [testing, setTesting] = useState(false);

  const handleFieldChange = (fieldName: string, value: string) => {
    setConfig((prev) => ({ ...prev, [fieldName]: value }));
    // Clear error for this field
    setErrors((prev) => {
      const newErrors = { ...prev };
      delete newErrors[fieldName];
      return newErrors;
    });
  };

  const validate = (): boolean => {
    const newErrors: Record<string, string> = {};

    providerType.config_schema.forEach((field) => {
      const value = config[field.name];
      if (field.required && (!value || value.trim() === '')) {
        newErrors[field.name] = `${field.display_name} is required`;
      }
    });

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleTest = async () => {
    if (!onTest) return;
    if (!validate()) return;

    setTesting(true);
    try {
      await onTest(config);
    } finally {
      setTesting(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) return;

    setSubmitting(true);
    try {
      const success = await onSubmit({
        display_name: displayName,
        enabled,
        config,
      });
      if (success) {
        // Form will be closed by parent
      }
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Provider Type Info */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <h3 className="font-semibold text-blue-900 mb-1">
          {providerType.display_name}
        </h3>
        <p className="text-sm text-blue-700">{providerType.description}</p>
        {providerType.requires_auth && (
          <p className="text-xs text-blue-600 mt-2">
            ⚠️ This provider requires authentication
          </p>
        )}
      </div>

      {/* Display Name */}
      <Input
        label="Display Name"
        value={displayName}
        onChange={(e) => setDisplayName(e.target.value)}
        required
        help="A friendly name for this provider configuration"
      />

      {/* Enabled Toggle */}
      <div className="flex items-center gap-3">
        <input
          type="checkbox"
          id="enabled"
          checked={enabled}
          onChange={(e) => setEnabled(e.target.checked)}
          className="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
        />
        <label htmlFor="enabled" className="text-sm font-medium text-gray-700">
          Enable this provider
        </label>
      </div>

      {/* Dynamic Config Fields */}
      <div className="space-y-4 pt-2">
        <h4 className="font-medium text-gray-900">Configuration</h4>
        {providerType.config_schema.map((field) => (
          <ProviderFormField
            key={field.name}
            field={field}
            value={config[field.name] || ''}
            onChange={(value) => handleFieldChange(field.name, value)}
            error={errors[field.name]}
          />
        ))}
      </div>

      {/* Action Buttons */}
      <div className="flex items-center justify-end gap-3 pt-4">
        {onTest && (
          <Button
            type="button"
            variant="secondary"
            onClick={handleTest}
            loading={testing}
            disabled={submitting}
          >
            Test Connection
          </Button>
        )}
        <Button
          type="button"
          variant="ghost"
          onClick={onCancel}
          disabled={submitting || testing}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          variant="primary"
          loading={submitting}
          disabled={testing}
        >
          {existingConfig ? 'Update Provider' : 'Add Provider'}
        </Button>
      </div>
    </form>
  );
};
