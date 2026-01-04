import React, { useState } from 'react';
import { Input } from '../common/Input';
import type { ProviderConfigField } from '../../types/provider';

interface ProviderFormFieldProps {
  field: ProviderConfigField;
  value: string;
  onChange: (value: string) => void;
  error?: string;
}

export const ProviderFormField: React.FC<ProviderFormFieldProps> = ({
  field,
  value,
  onChange,
  error,
}) => {
  const [showPassword, setShowPassword] = useState(false);

  // Render based on field type
  switch (field.type.toLowerCase()) {
    case 'secret':
    case 'password':
      return (
        <div className="relative">
          <Input
            type={showPassword ? 'text' : 'password'}
            label={field.display_name}
            value={value}
            onChange={(e) => onChange(e.target.value)}
            required={field.required}
            error={error}
            help={field.description}
            placeholder={field.default}
          />
          <button
            type="button"
            onClick={() => setShowPassword(!showPassword)}
            className="absolute right-3 top-9 text-gray-400 hover:text-gray-600"
          >
            {showPassword ? (
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"
                />
              </svg>
            ) : (
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                />
              </svg>
            )}
          </button>
        </div>
      );

    case 'url':
      return (
        <Input
          type="url"
          label={field.display_name}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          required={field.required}
          error={error}
          help={field.description}
          placeholder={field.default}
        />
      );

    case 'number':
      return (
        <Input
          type="number"
          label={field.display_name}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          required={field.required}
          error={error}
          help={field.description}
          placeholder={field.default}
        />
      );

    case 'email':
      return (
        <Input
          type="email"
          label={field.display_name}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          required={field.required}
          error={error}
          help={field.description}
          placeholder={field.default}
        />
      );

    case 'text':
    case 'string':
    default:
      return (
        <Input
          type="text"
          label={field.display_name}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          required={field.required}
          error={error}
          help={field.description}
          placeholder={field.default}
        />
      );
  }
};
