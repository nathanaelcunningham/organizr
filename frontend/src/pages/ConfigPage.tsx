import React, { useEffect } from 'react';
import { NavLink, Outlet } from 'react-router-dom';
import { PageHeader } from '../components/layout/PageHeader';
import { ProviderList } from '../components/providers/ProviderList';
import { ConfigForm } from '../components/config/ConfigForm';
import { useProviderStore } from '../stores/useProviderStore';
import { useConfigStore } from '../stores/useConfigStore';

export const ConfigPage: React.FC = () => {
  const { fetchProviders, fetchProviderTypes } = useProviderStore();
  const { fetchConfig } = useConfigStore();

  useEffect(() => {
    fetchProviders();
    fetchProviderTypes();
    fetchConfig();
  }, [fetchProviders, fetchProviderTypes, fetchConfig]);

  return (
    <div>
      <PageHeader
        title="Configuration"
        subtitle="Manage providers and application settings"
      />

      {/* Tab Navigation */}
      <div className="mb-6 border-b border-gray-200 -mx-4 sm:mx-0 px-4 sm:px-0">
        <nav className="-mb-px flex gap-4 sm:gap-6 overflow-x-auto scrollbar-hide">
          <NavLink
            to="/config/providers"
            className={({ isActive }) =>
              `py-2 px-1 border-b-2 font-medium text-sm transition-colors ${
                isActive
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`
            }
          >
            Providers
          </NavLink>
          <NavLink
            to="/config/general"
            className={({ isActive }) =>
              `py-2 px-1 border-b-2 font-medium text-sm transition-colors ${
                isActive
                  ? 'border-blue-600 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              }`
            }
          >
            General Settings
          </NavLink>
        </nav>
      </div>

      {/* Nested routes render here */}
      <Outlet />
    </div>
  );
};

// Provider Config Sub-page
export const ProvidersConfigPage: React.FC = () => {
  const { providers, loading } = useProviderStore();

  return (
    <div>
      <ProviderList providers={providers} loading={loading} />
    </div>
  );
};

// General Config Sub-page
export const GeneralConfigPage: React.FC = () => {
  return (
    <div>
      <ConfigForm />
    </div>
  );
};
