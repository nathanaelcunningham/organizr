import React, { useMemo } from 'react';
import { SearchResultListItem } from './SearchResultListItem';
import type { SearchResult } from '../../types/search';

export interface SeriesGroupProps {
    seriesName: string;
    books: SearchResult[];
    batchMode?: boolean;
    selectedIds?: Set<string>;
    onToggleSelection?: (result: SearchResult) => void;
    isSelected?: (result: SearchResult) => boolean;
}

/**
 * Displays a group of books from the same series with visual separation.
 * Shows series name with book count and renders books using SearchResultListItem.
 * Supports batch selection mode with indeterminate checkbox state.
 */
export const SeriesGroup: React.FC<SeriesGroupProps> = ({
    seriesName,
    books,
    batchMode = false,
    onToggleSelection,
    isSelected,
}) => {
    // Calculate selection state for series group checkbox
    const selectionState = useMemo(() => {
        if (!batchMode || !isSelected) return 'none';

        const selectedCount = books.filter(book => isSelected(book)).length;
        if (selectedCount === 0) return 'none';
        if (selectedCount === books.length) return 'all';
        return 'some';
    }, [books, batchMode, isSelected]);

    // Handle series group checkbox click
    const handleSeriesCheckboxClick = () => {
        if (!onToggleSelection) return;

        // If all selected, deselect all; otherwise, select all
        if (selectionState === 'all') {
            books.forEach(book => {
                if (isSelected && isSelected(book)) {
                    onToggleSelection(book);
                }
            });
        } else {
            books.forEach(book => {
                if (isSelected && !isSelected(book)) {
                    onToggleSelection(book);
                }
            });
        }
    };

    return (
        <div className="mb-6">
            {/* Series header with book count and optional checkbox */}
            <div className="flex items-center gap-3 px-4 py-2 bg-gray-100 rounded mb-3">
                {batchMode && (
                    <div onClick={(e) => e.stopPropagation()}>
                        <input
                            type="checkbox"
                            checked={selectionState === 'all'}
                            ref={(input) => {
                                if (input) {
                                    input.indeterminate = selectionState === 'some';
                                }
                            }}
                            onChange={handleSeriesCheckboxClick}
                            className="w-4 h-4 text-blue-600 rounded focus:ring-blue-500 cursor-pointer"
                        />
                    </div>
                )}
                <h3 className="text-lg font-semibold">
                    {seriesName}
                    <span className="text-sm font-normal text-gray-600 ml-2">
                        ({books.length} {books.length === 1 ? 'book' : 'books'})
                    </span>
                </h3>
            </div>

            {/* Books list with clean separation */}
            <div className="divide-y divide-gray-200 border border-gray-200 rounded-lg overflow-hidden">
                {books.map((book, idx) => (
                    <SearchResultListItem
                        key={`${book.id || book.title}-${idx}`}
                        result={book}
                        showSeriesNumber={seriesName !== 'Standalone'}
                        batchMode={batchMode}
                        selected={isSelected ? isSelected(book) : false}
                        onToggleSelect={onToggleSelection}
                    />
                ))}
            </div>
        </div>
    );
};
