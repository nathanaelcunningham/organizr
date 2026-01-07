import React from 'react';
import { SearchResultListItem } from './SearchResultListItem';
import type { SearchResult } from '../../types/search';

export interface SeriesGroupProps {
    seriesName: string;
    books: SearchResult[];
}

/**
 * Displays a group of books from the same series with visual separation.
 * Shows series name with book count and renders books using SearchResultListItem.
 */
export const SeriesGroup: React.FC<SeriesGroupProps> = ({ seriesName, books }) => {
    return (
        <div className="mb-6">
            {/* Series header with book count */}
            <h3 className="text-lg font-semibold mb-3 px-4 py-2 bg-gray-100 rounded">
                {seriesName}
                <span className="text-sm font-normal text-gray-600 ml-2">
                    ({books.length} {books.length === 1 ? 'book' : 'books'})
                </span>
            </h3>

            {/* Books list with clean separation */}
            <div className="divide-y divide-gray-200 border border-gray-200 rounded-lg overflow-hidden">
                {books.map((book, idx) => (
                    <SearchResultListItem
                        key={`${book.id || book.title}-${idx}`}
                        result={book}
                        showSeriesNumber={seriesName !== 'Standalone'}
                    />
                ))}
            </div>
        </div>
    );
};
