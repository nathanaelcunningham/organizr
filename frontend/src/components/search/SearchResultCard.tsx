import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card } from '../common/Card';
import { Button } from '../common/Button';
import { Badge } from '../common/Badge';
import type { SearchResult } from '../../types/search';
import { useDownloadStore } from '../../stores/useDownloadStore';
import { formatFileSize, truncate } from '../../utils/formatters';

interface SearchResultCardProps {
  result: SearchResult;
}

export const SearchResultCard: React.FC<SearchResultCardProps> = ({
  result,
}) => {
  const navigate = useNavigate();
  const createDownload = useDownloadStore((state) => state.createDownload);
  const [downloading, setDownloading] = useState(false);
  const [showDescription, setShowDescription] = useState(false);

  const handleDownload = async () => {
    setDownloading(true);
    try {
      // Convert series array to string representation
      const seriesString = result.series && result.series.length > 0
        ? result.series.map(s => `${s.name} #${s.number}`).join(', ')
        : undefined;

      await createDownload({
        title: result.title,
        author: result.author,
        series: seriesString,
        category: 'Audiobooks',
        torrent_url: result.torrent_url,
        magnet_link: result.magnet_link,
      });
      // Navigate to downloads page after successful creation
      navigate('/downloads');
    } catch (error) {
      console.error('Failed to create download:', error);
    } finally {
      setDownloading(false);
    }
  };

  return (
    <Card className="hover:shadow-md transition-shadow">
      <div className="space-y-3">
        {/* Title and Author */}
        <div>
          <h3 className="text-lg font-semibold text-gray-900">
            {result.title}
          </h3>
          <p className="text-sm text-gray-600 mt-1">by {result.author}</p>
          {result.series && result.series.length > 0 && (
            <p className="text-sm text-gray-500 mt-1">
              Series: {result.series.map(s => `${s.name} #${s.number}`).join(', ')}
            </p>
          )}
          {result.narrator && (
            <p className="text-sm text-gray-500">
              Narrated by: {result.narrator}
            </p>
          )}
        </div>

        {/* Metadata Grid */}
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3 text-sm">
          <div>
            <span className="text-gray-500">Size:</span>
            <span className="ml-2 font-medium text-gray-900">
              {formatFileSize(result.size)}
            </span>
          </div>
          <div>
            <span className="text-gray-500">Seeders:</span>
            <span
              className={`ml-2 font-medium ${
                result.seeders > 10 ? 'text-green-600' : 'text-yellow-600'
              }`}
            >
              {result.seeders}
            </span>
          </div>
          <div>
            <span className="text-gray-500">Leechers:</span>
            <span className="ml-2 font-medium text-gray-900">
              {result.leechers}
            </span>
          </div>
          <div>
            <span className="text-gray-500">Provider:</span>
            <span className="ml-2 font-medium text-gray-900">
              {result.provider}
            </span>
          </div>
        </div>

        {/* Badges */}
        <div className="flex flex-wrap gap-2">
          {result.freeleech && (
            <Badge variant="success" size="sm">
              Freeleech
            </Badge>
          )}
          {result.freeleech_vip && (
            <Badge variant="info" size="sm">
              VIP Freeleech
            </Badge>
          )}
          {result.vip && (
            <Badge variant="warning" size="sm">
              VIP
            </Badge>
          )}
          {result.file_type && (
            <Badge variant="default" size="sm">
              {result.file_type}
            </Badge>
          )}
          {result.language && (
            <Badge variant="default" size="sm">
              {result.language}
            </Badge>
          )}
          {result.category && (
            <Badge variant="default" size="sm">
              {result.category}
            </Badge>
          )}
        </div>

        {/* Tags */}
        {result.tags && result.tags.length > 0 && (
          <div className="flex flex-wrap gap-1">
            {result.tags.slice(0, 5).map((tag, index) => (
              <span
                key={index}
                className="inline-block px-2 py-1 text-xs bg-gray-100 text-gray-600 rounded"
              >
                {tag}
              </span>
            ))}
            {result.tags.length > 5 && (
              <span className="inline-block px-2 py-1 text-xs text-gray-500">
                +{result.tags.length - 5} more
              </span>
            )}
          </div>
        )}

        {/* Description */}
        {result.description && (
          <div>
            <button
              onClick={() => setShowDescription(!showDescription)}
              className="text-sm text-blue-600 hover:text-blue-700 font-medium"
            >
              {showDescription ? 'Hide' : 'Show'} Description
            </button>
            {showDescription && (
              <p className="mt-2 text-sm text-gray-700 whitespace-pre-wrap">
                {truncate(result.description, 500)}
              </p>
            )}
          </div>
        )}

        {/* Additional Info */}
        {(result.num_files || result.times_completed || result.added) && (
          <div className="flex flex-wrap gap-4 text-xs text-gray-500">
            {result.num_files && <span>Files: {result.num_files}</span>}
            {result.times_completed && (
              <span>Downloads: {result.times_completed}</span>
            )}
            {result.added && <span>Added: {result.added}</span>}
          </div>
        )}

        {/* Download Button */}
        <div className="pt-2">
          <Button
            variant="primary"
            size="md"
            onClick={handleDownload}
            loading={downloading}
            disabled={!result.torrent_url && !result.magnet_link}
            className="w-full"
          >
            {!result.torrent_url && !result.magnet_link
              ? 'No Download Available'
              : 'Download'}
          </Button>
        </div>
      </div>
    </Card>
  );
};
