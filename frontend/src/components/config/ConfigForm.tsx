import React, { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { Button } from '../common/Button';
import { Input } from '../common/Input';
import { Select } from '../common/Select';
import { ConfigSection } from './ConfigSection';
import { useConfigStore } from '../../stores/useConfigStore';
import { useNotificationStore } from '../../stores/useNotificationStore';
import { CONFIG_KEYS } from '../../types/config';
import { searchApi } from '../../api/search';

// Field configuration for cleaner form management
const FIELD_CONFIG = {
    qbittorrentUrl: { key: CONFIG_KEYS.QBITTORRENT_URL, default: 'http://localhost:8080' },
    qbittorrentUsername: { key: CONFIG_KEYS.QBITTORRENT_USERNAME, default: 'admin' },
    qbittorrentPassword: { key: CONFIG_KEYS.QBITTORRENT_PASSWORD, default: '' },
    pathsDestination: { key: CONFIG_KEYS.PATHS_DESTINATION, default: '/audiobooks' },
    pathsTemplate: { key: CONFIG_KEYS.PATHS_TEMPLATE, default: '{author}/{series}/{title}' },
    pathsNoSeriesTemplate: { key: CONFIG_KEYS.PATHS_NO_SERIES_TEMPLATE, default: '{author}/{title}' },
    pathsOperation: { key: CONFIG_KEYS.PATHS_OPERATION, default: 'copy' },
    monitorInterval: { key: CONFIG_KEYS.MONITOR_INTERVAL, default: '30' },
    monitorAutoOrganize: { key: CONFIG_KEYS.MONITOR_AUTO_ORGANIZE, default: 'true' },
    mamBaseUrl: { key: CONFIG_KEYS.MAM_BASEURL, default: 'https://www.myanonamouse.net' },
    mamSecret: { key: CONFIG_KEYS.MAM_SECRET, default: '' },
} as const;

type FormFieldKey = keyof typeof FIELD_CONFIG;
type FormData = Record<FormFieldKey, string>;

// Get default values from FIELD_CONFIG
const getDefaultValues = (): FormData => {
    return Object.entries(FIELD_CONFIG).reduce((acc, [field, { default: defaultValue }]) => {
        acc[field as FormFieldKey] = defaultValue;
        return acc;
    }, {} as FormData);
};

export function ConfigForm() {
    const { config, loading, updateMultipleConfigs } = useConfigStore();

    const { register, handleSubmit, reset, formState: { isSubmitting, dirtyFields } } = useForm<FormData>({
        defaultValues: getDefaultValues(),
    });
    console.log(config)

    // When config loads, populate the form
    useEffect(() => {
        if (config && Object.keys(config).length > 0) {
            const formData = Object.entries(FIELD_CONFIG).reduce((acc, [field, { key, default: defaultValue }]) => {
                console.log(config[key])
                acc[field as FormFieldKey] = config[key] ?? defaultValue;
                return acc;
            }, {} as FormData);

            console.log('Config loaded, resetting form:', formData);
            reset(formData);
        }
    }, [config, reset]);

    const onSubmit = async (data: FormData) => {
        try {
            // Build updates object with ONLY changed (dirty) fields
            const updates: Record<string, string> = {};

            Object.entries(FIELD_CONFIG).forEach(([field, { key }]) => {
                if (dirtyFields[field as FormFieldKey]) {
                    updates[key] = data[field as FormFieldKey];
                }
            });

            // Remove any empty values
            Object.keys(updates).forEach(key => {
                if (updates[key] === '') {
                    delete updates[key];
                }
            });

            // Only submit if there are changes
            if (Object.keys(updates).length === 0) {
                useNotificationStore.getState().addNotification('info', 'No changes to save');
                return;
            }

            console.log('Submitting updates:', updates);
            const success = await updateMultipleConfigs(updates);

            // Reset form with current values to clear dirty state
            if (success) {
                reset(data);
            }
        } catch (error) {
            console.error('Form submission error:', error);
        }
    };

    const testConnection = async () => {
        const res = await searchApi.testConnection()
        if (res.success) {
            useNotificationStore.getState().addNotification('info', 'Succesfull connection')
        } else {
            useNotificationStore.getState().addNotification('error', 'Failed connection')
        }

    }


    if (loading) {
        return <div className="text-gray-500">Loading configuration...</div>;
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
                <Button type="button" onClick={testConnection}>Test Connection</Button>
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
                <Input
                    label="Username"
                    type="text"
                    {...register('qbittorrentUsername')}
                    required
                />
                <Input
                    label="Password"
                    type="password"
                    {...register('qbittorrentPassword')}
                    help="Leave blank to keep existing password"
                />
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
                <Input
                    label="Path Template (with series)"
                    type="text"
                    {...register('pathsTemplate')}
                    required
                    help="Template for organizing files with series. Variables: {author}, {series}, {title}"
                />
                <Input
                    label="Path Template (without series)"
                    type="text"
                    {...register('pathsNoSeriesTemplate')}
                    required
                    help="Template for organizing files without series. Variables: {author}, {title}"
                />
                <Select
                    label="Operation"
                    {...register('pathsOperation')}
                    options={[
                        { value: 'copy', label: 'Copy files' },
                        { value: 'move', label: 'Move files' },
                    ]}
                    help="Whether to copy or move files to the organized location"
                />
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
                <div className="flex items-center gap-3">
                    <input
                        type="checkbox"
                        id="autoOrganize"
                        {...register('monitorAutoOrganize')}
                        value="true"
                        className="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
                    />
                    <label
                        htmlFor="autoOrganize"
                        className="text-sm font-medium text-gray-700"
                    >
                        Automatically organize completed downloads
                    </label>
                </div>
            </ConfigSection>

            {/* Submit Button */}
            <div className="flex justify-end">
                <Button type="submit" variant="primary" loading={isSubmitting}>
                    Save Settings
                </Button>
            </div>
        </form>
    );
};
