import { useEffect, useState } from 'react'
import { useForm, useWatch } from 'react-hook-form'
import { Button } from '../common/Button'
import { Input } from '../common/Input'
import { Select } from '../common/Select'
import { Spinner } from '../common/Spinner'
import { ConfigSection } from './ConfigSection'
import { useConfigStore } from '../../stores/useConfigStore'
import { useNotificationStore } from '../../stores/useNotificationStore'
import { CONFIG_KEYS } from '../../types/config'
import { searchApi } from '../../api/search'
import { qbittorrentApi } from '../../api/qbittorrent'
import { configApi, type PreviewPathResponse } from '../../api/config'

// Field configuration for cleaner form management
const FIELD_CONFIG = {
  qbittorrentUrl: { key: CONFIG_KEYS.QBITTORRENT_URL, default: 'http://localhost:8080' },
  qbittorrentUsername: { key: CONFIG_KEYS.QBITTORRENT_USERNAME, default: 'admin' },
  qbittorrentPassword: { key: CONFIG_KEYS.QBITTORRENT_PASSWORD, default: '' },
  pathsDestination: { key: CONFIG_KEYS.PATHS_DESTINATION, default: '/audiobooks' },
  pathsTemplate: { key: CONFIG_KEYS.PATHS_TEMPLATE, default: '{author}/{series}/{title}' },
  pathsNoSeriesTemplate: { key: CONFIG_KEYS.PATHS_NO_SERIES_TEMPLATE, default: '{author}/{title}' },
  pathsOperation: { key: CONFIG_KEYS.PATHS_OPERATION, default: 'copy' },
  pathsLocalMount: { key: CONFIG_KEYS.PATHS_LOCAL_MOUNT, default: '' },
  monitorInterval: { key: CONFIG_KEYS.MONITOR_INTERVAL, default: '30' },
  organizationAutoOrganize: { key: CONFIG_KEYS.ORGANIZATION_AUTO_ORGANIZE, default: 'true' },
  mamBaseUrl: { key: CONFIG_KEYS.MAM_BASEURL, default: 'https://www.myanonamouse.net' },
  mamSecret: { key: CONFIG_KEYS.MAM_SECRET, default: '' },
} as const

type FormFieldKey = keyof typeof FIELD_CONFIG
type FormData = Record<FormFieldKey, string>

// Get default values from FIELD_CONFIG
const getDefaultValues = (): FormData => {
  return Object.entries(FIELD_CONFIG).reduce((acc, [field, { default: defaultValue }]) => {
    acc[field as FormFieldKey] = defaultValue
    return acc
  }, {} as FormData)
}

// PathPreview component for real-time template preview
function PathPreview({ template, hasSeries }: { template: string; hasSeries: boolean }) {
  const [preview, setPreview] = useState<PreviewPathResponse | null>(null)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    if (!template) {
      setPreview(null)
      return
    }

    // Debounce the preview API call
    const timeoutId = setTimeout(async () => {
      setLoading(true)
      try {
        const result = await configApi.previewPath({
          template,
          author: 'Example Author',
          series: hasSeries ? 'Example Series' : undefined,
          series_number: hasSeries ? '1' : undefined,
          title: 'Example Book Title',
        })
        setPreview(result)
      } catch (error) {
        setPreview({
          valid: false,
          error: error instanceof Error ? error.message : 'Failed to preview path',
        })
      } finally {
        setLoading(false)
      }
    }, 500)

    return () => clearTimeout(timeoutId)
  }, [template, hasSeries])

  if (!template) return null

  if (loading) {
    return <div className="text-xs text-gray-500 mt-1">Generating preview...</div>
  }

  if (!preview) return null

  if (preview.valid && preview.path) {
    return (
      <div className="text-xs text-gray-600 mt-1">
        <span className="font-medium">Preview:</span>{' '}
        <span className="font-mono">/audiobooks/{preview.path}</span>
      </div>
    )
  }

  if (preview.error) {
    return <div className="text-xs text-red-600 mt-1">{preview.error}</div>
  }

  return null
}

export function ConfigForm() {
  const { config, loading, updateMultipleConfigs } = useConfigStore()
  const [qbTestLoading, setQbTestLoading] = useState(false)
  const [qbTestResult, setQbTestResult] = useState<{ success: boolean; message: string } | null>(
    null
  )

  const {
    register,
    handleSubmit,
    reset,
    control,
    formState: { isSubmitting, dirtyFields },
  } = useForm<FormData>({
    defaultValues: getDefaultValues(),
  })

  // Watch template fields for preview
  const pathsTemplate = useWatch({ control, name: 'pathsTemplate' })
  const pathsNoSeriesTemplate = useWatch({ control, name: 'pathsNoSeriesTemplate' })

  console.log(config)

  // When config loads, populate the form
  useEffect(() => {
    if (config && Object.keys(config).length > 0) {
      const formData = Object.entries(FIELD_CONFIG).reduce(
        (acc, [field, { key, default: defaultValue }]) => {
          console.log(config[key])
          acc[field as FormFieldKey] = config[key] ?? defaultValue
          return acc
        },
        {} as FormData
      )

      console.log('Config loaded, resetting form:', formData)
      reset(formData)
    }
  }, [config, reset])

  const onSubmit = async (data: FormData) => {
    try {
      // Build updates object with ONLY changed (dirty) fields
      const updates: Record<string, string> = {}

      Object.entries(FIELD_CONFIG).forEach(([field, { key }]) => {
        if (dirtyFields[field as FormFieldKey]) {
          let value = data[field as FormFieldKey]

          // Handle checkbox fields - convert boolean-like values to string
          if (field === 'organizationAutoOrganize') {
            // Checkbox can be empty string (unchecked) or "on"/"true" (checked)
            value = value ? 'true' : 'false'
          }

          updates[key] = value
        }
      })

      // Remove any empty values (except 'false' for booleans)
      Object.keys(updates).forEach((key) => {
        if (updates[key] === '' && key !== 'organization.auto_organize') {
          delete updates[key]
        }
      })

      // Only submit if there are changes
      if (Object.keys(updates).length === 0) {
        useNotificationStore.getState().addNotification('info', 'No changes to save')
        return
      }

      console.log('Submitting updates:', updates)
      const success = await updateMultipleConfigs(updates)

      // Reset form with current values to clear dirty state
      if (success) {
        reset(data)
      }
    } catch (error) {
      console.error('Form submission error:', error)
    }
  }

  const testConnection = async () => {
    const res = await searchApi.testConnection()
    if (res.success) {
      useNotificationStore.getState().addNotification('info', 'Succesfull connection')
    } else {
      useNotificationStore.getState().addNotification('error', 'Failed connection')
    }
  }

  const testQBittorrentConnection = async () => {
    setQbTestResult(null)
    setQbTestLoading(true)
    try {
      const res = await qbittorrentApi.testConnection()
      setQbTestResult(res)
    } catch (error) {
      setQbTestResult({
        success: false,
        message: error instanceof Error ? error.message : 'Unknown error',
      })
    } finally {
      setQbTestLoading(false)
    }
  }

  if (loading) {
    return <div className="text-gray-500">Loading configuration...</div>
  }

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      {/* MyAnonamouse Configuration */}
      <ConfigSection
        title="MyAnonamouse Search"
        description="Configure MyAnonamouse (MAM) API credentials for searching audiobooks"
      >
        <Input
          label="Base URL"
          type="url"
          {...register('mamBaseUrl')}
          required
          help="MyAnonamouse base URL (default: https://www.myanonamouse.net)"
        />
        <Input
          label="API Secret"
          type="text"
          {...register('mamSecret')}
          required
          help="Your MyAnonamouse API secret/key (find this in your MAM account settings)"
        />
        <Button type="button" onClick={testConnection}>
          Test Connection
        </Button>
      </ConfigSection>

      {/* qBittorrent Connection */}
      <ConfigSection
        title="qBittorrent Connection"
        description="Configure connection to your qBittorrent instance"
      >
        <Input
          label="Web UI URL"
          type="url"
          {...register('qbittorrentUrl')}
          required
          help="URL to qBittorrent Web UI (e.g., http://localhost:8080)"
        />
        <Input label="Username" type="text" {...register('qbittorrentUsername')} required />
        <Input
          label="Password"
          type="password"
          {...register('qbittorrentPassword')}
          help="Leave blank to keep existing password"
        />
        <div className="space-y-2">
          <Button
            type="button"
            variant="secondary"
            onClick={testQBittorrentConnection}
            disabled={qbTestLoading}
          >
            {qbTestLoading ? 'Testing...' : 'Test Connection'}
          </Button>
          {qbTestLoading && (
            <div className="flex items-center gap-2 text-sm text-gray-600">
              <Spinner size="sm" />
              <span>Testing connection...</span>
            </div>
          )}
          {qbTestResult && !qbTestLoading && (
            <div
              className={`flex items-center gap-2 text-sm ${qbTestResult.success ? 'text-green-600' : 'text-red-600'}`}
            >
              <span>{qbTestResult.success ? '✓' : '✗'}</span>
              <span>{qbTestResult.message}</span>
            </div>
          )}
        </div>
      </ConfigSection>

      {/* File Organization */}
      <ConfigSection
        title="File Organization"
        description="Configure how downloaded audiobooks are organized"
      >
        <Input
          label="Destination Directory"
          type="text"
          {...register('pathsDestination')}
          required
          help="Base directory where audiobooks will be organized"
        />
        <div>
          <Input
            label="Path Template (with series)"
            type="text"
            {...register('pathsTemplate')}
            required
            help="Template for organizing files with series. Variables: {author}, {series}, {series_number}, {title}"
          />
          <PathPreview template={pathsTemplate || ''} hasSeries={true} />
        </div>
        <div>
          <Input
            label="Path Template (without series)"
            type="text"
            {...register('pathsNoSeriesTemplate')}
            required
            help="Template for organizing files without series. Variables: {author}, {title}"
          />
          <PathPreview template={pathsNoSeriesTemplate || ''} hasSeries={false} />
        </div>
        <Select
          label="Operation"
          {...register('pathsOperation')}
          options={[
            { value: 'copy', label: 'Copy files' },
            { value: 'move', label: 'Move files' },
          ]}
          help="Whether to copy or move files to the organized location"
        />
        <div className="pt-4 border-t border-gray-200">
          <h4 className="text-sm font-medium text-gray-700 mb-2">Remote qBittorrent Setup</h4>
          <p className="text-xs text-gray-500 mb-4">
            For network shares or Docker deployments: configure where qBittorrent's files are
            accessible
          </p>
          <Input
            label="Download Mount Point"
            type="text"
            {...register('pathsLocalMount')}
            placeholder="/Volumes/data"
            help="Where qBittorrent's download directory is mounted locally (e.g., /Volumes/data on macOS, or Docker volume mount). Leave empty if running locally."
          />
          <p className="text-xs text-gray-500 mt-2">
            Example: If qBittorrent saves to /downloads and you mount the parent at /Volumes/data,
            enter /Volumes/data
          </p>
        </div>
      </ConfigSection>

      {/* Monitor Settings */}
      <ConfigSection
        title="Monitor Settings"
        description="Configure the background monitor that checks for completed downloads"
      >
        <Input
          label="Interval (seconds)"
          type="number"
          {...register('monitorInterval')}
          required
          help="How often to check for completed downloads"
        />
        <div className="space-y-2">
          <div className="flex items-center gap-3">
            <input
              type="checkbox"
              id="autoOrganize"
              {...register('organizationAutoOrganize')}
              className="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
            />
            <label htmlFor="autoOrganize" className="text-sm font-medium text-gray-700">
              Automatically organize completed downloads
            </label>
          </div>
          <p className="text-xs text-gray-500 ml-7">
            When enabled, downloads will be automatically organized after completion. Disable to
            organize manually.
          </p>
        </div>
      </ConfigSection>

      {/* Submit Button */}
      <div className="flex justify-end">
        <Button type="submit" variant="primary" loading={isSubmitting}>
          Save Settings
        </Button>
      </div>
    </form>
  )
}
