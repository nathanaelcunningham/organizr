import type { SearchResult } from '../types/search';

export interface SeriesGroup {
    seriesName: string;
    books: SearchResult[];
}

/**
 * Groups search results by series name and sorts books within each series by number.
 * Books without series are grouped under "Standalone".
 * Books with multiple series appear in each respective group.
 */
export function groupBySeries(results: SearchResult[]): SeriesGroup[] {
    // Group by series name using reduce
    const grouped = results.reduce((acc, result) => {
        if (!result.series || result.series.length === 0) {
            // Standalone books
            acc['Standalone'] = [...(acc['Standalone'] || []), result];
        } else {
            // Books with series - handle multiple series
            result.series.forEach(s => {
                acc[s.name] = [...(acc[s.name] || []), result];
            });
        }
        return acc;
    }, {} as Record<string, SearchResult[]>);

    // Convert to array and sort books within each series
    return Object.entries(grouped).map(([seriesName, books]) => ({
        seriesName,
        books: [...books].sort((a, b) => {
            // Get first series' number for each book
            const numA = parseFloat(a.series?.[0]?.number || '999');
            const numB = parseFloat(b.series?.[0]?.number || '999');
            return numA - numB;
        })
    }));
}
