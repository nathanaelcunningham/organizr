import React, { useState, useEffect } from 'react';
import { Button } from '../common/Button';
import { Input } from '../common/Input';
import { Select } from '../common/Select';
import { ConfigSection } from './ConfigSection';
import { useConfigStore } from '../../stores/useConfigStore';
import { CONFIG_KEYS } from '../../types/config';

export const ConfigForm: React.FC = () => {
  const { config, loading, updateMultipleConfigs } = useConfigStore();
  const [formData, setFormData] = useState({
    qbittorrentUrl: '',
    qbittorrentUsername: '',
    qbittorrentPassword: '',
    pathsDestination: '',
    pathsTemplate: '',
    pathsNoSeriesTemplate: '',
    pathsOperation: 'copy',
    monitorInterval: '',
    monitorAutoOrganize: 'true',
  });
  const [submitting, setSubmitting] = useState(false);

  // Initialize form with config values
  useEffect(() => {
    if (config) {
      setFormData({
        qbittorrentUrl:
          config[CONFIG_KEYS.QBITTORRENT_URL] || 'http://localhost:8080',
        qbittorrentUsername: config[CONFIG_KEYS.QBITTORRENT_USERNAME] || 'admin',
        qbittorrentPassword: config[CONFIG_KEYS.QBITTORRENT_PASSWORD] || '',
        pathsDestination: config[CONFIG_KEYS.PATHS_DESTINATION] || '/audiobooks',
        pathsTemplate:
          config[CONFIG_KEYS.PATHS_TEMPLATE] ||
          '{author}/{series}/{title}',
        pathsNoSeriesTemplate:
          config[CONFIG_KEYS.PATHS_NO_SERIES_TEMPLATE] ||
          '{author}/{title}',
        pathsOperation: config[CONFIG_KEYS.PATHS_OPERATION] || 'copy',
        monitorInterval: config[CONFIG_KEYS.MONITOR_INTERVAL] || '30',
        monitorAutoOrganize:
          config[CONFIG_KEYS.MONITOR_AUTO_ORGANIZE] || 'true',
      });
    }
  }, [config]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);

    try {
      const updates: Record<string, string> = {
        [CONFIG_KEYS.QBITTORRENT_URL]: formData.qbittorrentUrl,
        [CONFIG_KEYS.QBITTORRENT_USERNAME]: formData.qbittorrentUsername,
        [CONFIG_KEYS.QBITTORRENT_PASSWORD]: formData.qbittorrentPassword,
        [CONFIG_KEYS.PATHS_DESTINATION]: formData.pathsDestination,
        [CONFIG_KEYS.PATHS_TEMPLATE]: formData.pathsTemplate,
        [CONFIG_KEYS.PATHS_NO_SERIES_TEMPLATE]:
          formData.pathsNoSeriesTemplate,
        [CONFIG_KEYS.PATHS_OPERATION]: formData.pathsOperation,
        [CONFIG_KEYS.MONITOR_INTERVAL]: formData.monitorInterval,
        [CONFIG_KEYS.MONITOR_AUTO_ORGANIZE]: formData.monitorAutoOrganize,
      };

      await updateMultipleConfigs(updates);
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return <div className="text-gray-500">Loading configuration...</div>;
  }

  return (
    <form onSubmit={handleSubmit}>
      {/* qBittorrent Connection */}
      <ConfigSection
        title="qBittorrent Connection"
        description="Configure connection to your qBittorrent instance"
      >
        <Input
          label="Web UI URL"
          type="url"
          value={formData.qbittorrentUrl}
          onChange={(e) =>
            setFormData({ ...formData, qbittorrentUrl: e.target.value })
          }
          required
          help="URL to qBittorrent Web UI (e.g., http://localhost:8080)"
        />
        <Input
          label="Username"
          type="text"
          value={formData.qbittorrentUsername}
          onChange={(e) =>
            setFormData({ ...formData, qbittorrentUsername: e.target.value })
          }
          required
        />
        <Input
          label="Password"
          type="password"
          value={formData.qbittorrentPassword}
          onChange={(e) =>
            setFormData({ ...formData, qbittorrentPassword: e.target.value })
          }
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
          value={formData.pathsDestination}
          onChange={(e) =>
            setFormData({ ...formData, pathsDestination: e.target.value })
          }
          required
          help="Base directory where audiobooks will be organized"
        />
        <Input
          label="Path Template (with series)"
          type="text"
          value={formData.pathsTemplate}
          onChange={(e) =>
            setFormData({ ...formData, pathsTemplate: e.target.value })
          }
          required
          help="Template for organizing files with series. Variables: {author}, {series}, {title}"
        />
        <Input
          label="Path Template (without series)"
          type="text"
          value={formData.pathsNoSeriesTemplate}
          onChange={(e) =>
            setFormData({
              ...formData,
              pathsNoSeriesTemplate: e.target.value,
            })
          }
          required
          help="Template for organizing files without series. Variables: {author}, {title}"
        />
        <Select
          label="Operation"
          value={formData.pathsOperation}
          onChange={(e) =>
            setFormData({ ...formData, pathsOperation: e.target.value })
          }
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
          value={formData.monitorInterval}
          onChange={(e) =>
            setFormData({ ...formData, monitorInterval: e.target.value })
          }
          required
          help="How often to check for completed downloads"
        />
        <div className="flex items-center gap-3">
          <input
            type="checkbox"
            id="autoOrganize"
            checked={formData.monitorAutoOrganize === 'true'}
            onChange={(e) =>
              setFormData({
                ...formData,
                monitorAutoOrganize: e.target.checked ? 'true' : 'false',
              })
            }
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
        <Button type="submit" variant="primary" loading={submitting}>
          Save Settings
        </Button>
      </div>
    </form>
  );
};
