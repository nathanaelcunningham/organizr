import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Button } from '../common/Button';
import { Badge } from '../common/Badge';
import type { SearchResult } from '../../types/search';
import { useDownloadStore } from '../../stores/useDownloadStore';
import { formatFileSize } from '../../utils/formatters';

interface SearchResultListItemProps {
    result: SearchResult;
    showSeriesNumber?: boolean;
}

export const SearchResultListItem: React.FC<SearchResultListItemProps> = ({
    result,
    showSeriesNumber = false,
}) => {
    const navigate = useNavigate();
    const createDownload = useDownloadStore((state) => state.createDownload);
    const [expanded, setExpanded] = useState(false);
    const [downloading, setDownloading] = useState(false);
    const [showDescription, setShowDescription] = useState(false);

    const handleDownload = async () => {
        setDownloading(true);
        try {
            // Convert series array to string representation (names only, not numbers)
            const seriesString = result.series && result.series.length > 0
                ? result.series.map(s => s.name).join(', ')
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
        <div className="border-b border-gray-200 last:border-b-0">
            {/* Compact Row */}
            <div
                className="flex flex-col md:flex-row md:items-center gap-3 px-4 py-3 hover:bg-gray-50 transition-colors cursor-pointer"
                onClick={() => setExpanded(!expanded)}
            >
                {/* Title & Author - Left side, flex-grow */}
                <div className="grow min-w-0 flex items-start justify-between md:block">
                    <div className="min-w-0 grow">
                        <h3 className="text-sm font-semibold text-gray-900 truncate">
                            {showSeriesNumber && result.series && result.series.length > 0 && (
                                <span className="text-gray-500 font-normal mr-2">
                                    #{result.series[0].number}
                                </span>
                            )}
                            {result.title}
                        </h3>
                        <p className="text-xs text-gray-600 truncate mt-0.5">
                            by {result.author}
                        </p>
                    </div>
                    {/* Seeders on mobile - right aligned */}
                    <div className="shrink-0 md:hidden ml-3">
                        <span className="text-xs text-gray-500">Seeders: </span>
                        <span
                            className={`text-sm font-medium ${result.seeders > 10 ? 'text-green-600' : 'text-yellow-600'
                                }`}
                        >
                            {result.seeders}
                        </span>
                    </div>
                </div>

                {/* Metadata columns - Desktop only */}
                <div className="hidden md:flex md:items-center md:gap-3 md:shrink-0">
                    {/* Seeders */}
                    <div className="w-16 text-right">
                        <span
                            className={`text-sm font-medium ${result.seeders > 10 ? 'text-green-600' : 'text-yellow-600'
                                }`}
                        >
                            {result.seeders}
                        </span>
                    </div>

                    {/* Size */}
                    <div className="w-20 text-right">
                        <span className="text-sm font-medium text-gray-900">
                            {formatFileSize(result.size)}
                        </span>
                    </div>


                    {/* Critical Badges */}
                    <div className="flex gap-1 shrink-0">
                        {result.freeleech && (
                            <Badge variant="success" size="sm">
                                FL
                            </Badge>
                        )}
                        {result.freeleech_vip && (
                            <Badge variant="info" size="sm">
                                VIP FL
                            </Badge>
                        )}
                    </div>

                    {/* Download Button */}
                    <div className="shrink-0" onClick={(e) => e.stopPropagation()}>
                        <Button
                            variant="primary"
                            size="sm"
                            onClick={handleDownload}
                            loading={downloading}
                            disabled={!result.torrent_url && !result.magnet_link}
                        >
                            {!result.torrent_url && !result.magnet_link ? 'N/A' : 'Download'}
                        </Button>
                    </div>

                    {/* Expand Toggle */}
                    <div className="shrink-0">
                        <svg
                            className={`w-5 h-5 text-gray-400 transition-transform ${expanded ? 'rotate-180' : ''
                                }`}
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M19 9l-7 7-7-7"
                            />
                        </svg>
                    </div>
                </div>

                {/* Mobile metadata row */}
                <div className="flex md:hidden items-center justify-between gap-3">
                    <div className="flex items-center gap-3 flex-wrap">
                        <span className="text-xs text-gray-600">
                            Size: <span className="font-medium text-gray-900">{formatFileSize(result.size)}</span>
                        </span>
                        {result.freeleech && (
                            <Badge variant="success" size="sm">
                                Freeleech
                            </Badge>
                        )}
                        {result.freeleech_vip && (
                            <Badge variant="info" size="sm">
                                VIP FL
                            </Badge>
                        )}
                    </div>
                    {/* Expand Toggle - Mobile */}
                    <div className="shrink-0">
                        <svg
                            className={`w-5 h-5 text-gray-400 transition-transform ${expanded ? 'rotate-180' : ''
                                }`}
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                        >
                            <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M19 9l-7 7-7-7"
                            />
                        </svg>
                    </div>
                </div>

                {/* Mobile download button */}
                <div className="md:hidden" onClick={(e) => e.stopPropagation()}>
                    <Button
                        variant="primary"
                        size="sm"
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

            {/* Expandable Section */}
            {expanded && (
                <div className="px-4 py-3 bg-gray-50 border-t border-gray-200">
                    <div className="space-y-3">
                        {/* Series & Narrator */}
                        {((result.series && result.series.length > 0) || result.narrator) && (
                            <div className="text-sm">
                                {result.series && result.series.length > 0 && (
                                    <p className="text-gray-700">
                                        <span className="font-medium">Series:</span> {result.series.map(s => `${s.name} #${s.number}`).join(', ')}
                                    </p>
                                )}
                                {result.narrator && (
                                    <p className="text-gray-700">
                                        <span className="font-medium">Narrated by:</span>{' '}
                                        {result.narrator}
                                    </p>
                                )}
                            </div>
                        )}

                        {/* Extended Metadata Grid */}
                        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 text-sm">
                            <div>
                                <span className="text-gray-500">Leechers:</span>
                                <span className="ml-2 font-medium text-gray-900">
                                    {result.leechers}
                                </span>
                            </div>
                            {result.num_files && (
                                <div>
                                    <span className="text-gray-500">Files:</span>
                                    <span className="ml-2 font-medium text-gray-900">
                                        {result.num_files}
                                    </span>
                                </div>
                            )}
                            {result.times_completed && (
                                <div>
                                    <span className="text-gray-500">Downloads:</span>
                                    <span className="ml-2 font-medium text-gray-900">
                                        {result.times_completed}
                                    </span>
                                </div>
                            )}
                            {result.added && (
                                <div>
                                    <span className="text-gray-500">Added:</span>
                                    <span className="ml-2 font-medium text-gray-900">
                                        {result.added}
                                    </span>
                                </div>
                            )}
                        </div>

                        {/* Additional Badges */}
                        <div className="flex flex-wrap gap-2">
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
                            <div>
                                <p className="text-xs font-medium text-gray-500 mb-1">Tags:</p>
                                <div className="flex flex-wrap gap-1">
                                    {result.tags.map((tag, index) => (
                                        <span
                                            key={index}
                                            className="inline-block px-2 py-1 text-xs bg-gray-100 text-gray-600 rounded"
                                        >
                                            {tag}
                                        </span>
                                    ))}
                                </div>
                            </div>
                        )}

                        {/* Description */}
                        {result.description && (
                            <div>
                                <button
                                    onClick={(e) => {
                                        e.stopPropagation();
                                        setShowDescription(!showDescription);
                                    }}
                                    className="text-sm text-blue-600 hover:text-blue-700 font-medium"
                                >
                                    {showDescription ? 'Hide' : 'Show'} Description
                                </button>
                                {showDescription && (
                                    <p className="mt-2 text-sm text-gray-700 whitespace-pre-wrap">
                                        {result.description.length > 500
                                            ? result.description.substring(0, 500) + '...'
                                            : result.description}
                                    </p>
                                )}
                            </div>
                        )}
                    </div>
                </div>
            )}
        </div>
    );
};
