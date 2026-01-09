import React, { useEffect } from 'react'
import { PageHeader } from '../components/layout/PageHeader'
import { ConfigForm } from '../components/config/ConfigForm'
import { useConfigStore } from '../stores/useConfigStore'

export const ConfigPage: React.FC = () => {
  const { fetchConfig } = useConfigStore()

  useEffect(() => {
    fetchConfig()
  }, [fetchConfig])

  return (
    <div>
      <PageHeader title="Configuration" subtitle="Manage application settings" />
      <ConfigForm />
    </div>
  )
}
