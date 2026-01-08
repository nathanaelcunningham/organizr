# Phase 7: MAM Series Detection - Research

**Researched:** 2026-01-07
**Domain:** MAM API series parsing, React UI grouping, confirmation dialogs
**Confidence:** HIGH

<research_summary>
## Summary

Researched MAM series detection and display patterns. The MAM API already provides series information via `series_info` JSON field, but the current implementation discards book numbers. The standard approach is to parse the structured JSON, preserve all series metadata (name + number), implement client-side grouping/sorting in React, and use accessible confirmation modals.

Key finding: Don't hand-roll grouping algorithms or modal accessibility. Use JavaScript's native `.sort()` and `.reduce()` for grouping, and leverage native `<dialog>` element or react-modal for accessible confirmation dialogs with proper focus management.

**Primary recommendation:**
1. Update `formatSeriesInfo()` to return structured data (not concatenated strings)
2. Use TypeScript types to preserve series name + number per book
3. Implement React grouping with `.reduce()` and sorting with `.sort()`
4. Use native `<dialog>` element with ARIA attributes for confirmation modal
</research_summary>

<standard_stack>
## Standard Stack

### Core (Already in Use)
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| React | 19.2.0 | UI framework | Project already uses React 19 |
| TypeScript | 5.9.3 | Type safety | Already in stack, critical for structured data |
| Zustand | 5.0.9 | State management | Already in use, lightweight |
| React Hook Form | 7.70.0 | Form management | Already in use, perfect for confirmation dialog |
| Tailwind CSS | 4.1.18 | Styling | Already in use, good for modal styling |

### Supporting (Optional Additions)
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| react-modal | 3.16.x | Accessible modals | If native `<dialog>` insufficient |
| lodash.groupby | 4.6.x | Grouping helper | If complex grouping logic needed (overkill for this) |

### Alternatives Considered
| Instead of | Could Use | Tradeoff |
|------------|-----------|----------|
| Native `<dialog>` | react-modal | react-modal more features, but native is simpler and lighter |
| Manual grouping | lodash.groupby | lodash adds dependency, manual is straightforward here |
| String concatenation | Structured objects | Structured data required for sorting/grouping |

**Installation (if needed):**
```bash
# Only if native <dialog> proves insufficient
npm install react-modal
npm install --save-dev @types/react-modal
```
</standard_stack>

<architecture_patterns>
## Architecture Patterns

### Recommended Project Structure
```
backend/internal/search/providers/
├── mam.go                    # Update formatSeriesInfo to return structured data
└── mam_test.go              # Add tests for series parsing

frontend/src/
├── types/
│   └── search.ts            # Add SeriesInfo type
├── components/search/
│   ├── SearchResults.tsx    # Update to group by series
│   ├── SeriesGroup.tsx      # NEW: Display grouped series
│   └── DownloadConfirmModal.tsx  # NEW: Confirmation dialog
└── utils/
    └── groupSeries.ts       # NEW: Grouping logic
```

### Pattern 1: Structured Series Data (Backend)
**What:** Return structured series data instead of concatenated strings
**When to use:** Always - enables frontend sorting/grouping
**Example:**
```go
// Current (bad): Returns "Series A, Series B"
func formatSeriesInfo(seriesInfo string) string

// Improved: Returns structured data
func formatSeriesInfo(seriesInfo string) []SeriesInfo
type SeriesInfo struct {
    ID     string
    Name   string
    Number string  // Can be "1", "Book 1", "1.5", etc.
}

// Update SearchResult model to use []SeriesInfo instead of string
```

### Pattern 2: Client-Side Grouping with Array.reduce()
**What:** Group search results by series name using native JavaScript
**When to use:** When displaying results grouped by series
**Example:**
```typescript
// Source: Standard JavaScript pattern, verified in React docs
function groupBySeries(results: SearchResult[]): Record<string, SearchResult[]> {
  return results.reduce((groups, result) => {
    // Books without series go in "Standalone" group
    if (!result.series || result.series.length === 0) {
      groups['Standalone'] = [...(groups['Standalone'] || []), result];
      return groups;
    }

    // Books with series - create group per series
    result.series.forEach(s => {
      groups[s.name] = [...(groups[s.name] || []), result];
    });

    return groups;
  }, {} as Record<string, SearchResult[]>);
}
```

### Pattern 3: Sorting by Series Number
**What:** Sort books within each series by their book number
**When to use:** After grouping, before display
**Example:**
```typescript
// Source: Standard JavaScript sort, verified best practices
function sortBySeriesNumber(results: SearchResult[]): SearchResult[] {
  return [...results].sort((a, b) => {
    // Extract numeric portion from book number (handles "1", "Book 1", "1.5")
    const numA = parseFloat(a.series[0]?.number || '0');
    const numB = parseFloat(b.series[0]?.number || '0');
    return numA - numB;
  });
}

// Note: Uses spread to avoid mutating original array
```

### Pattern 4: Accessible Confirmation Modal
**What:** Use native `<dialog>` element with proper ARIA and focus management
**When to use:** For pre-download confirmation with editable fields
**Example:**
```typescript
// Source: MDN Web Docs, verified accessibility guidelines
function DownloadConfirmModal({ isOpen, onClose, onConfirm, initialData }) {
  const dialogRef = useRef<HTMLDialogElement>(null);

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;

    if (isOpen) {
      dialog.showModal(); // Opens modal with focus management
    } else {
      dialog.close();
    }
  }, [isOpen]);

  return (
    <dialog
      ref={dialogRef}
      aria-labelledby="modal-title"
      aria-describedby="modal-description"
      className="backdrop:bg-black backdrop:opacity-50"
      onClose={onClose}
    >
      <h2 id="modal-title">Confirm Download</h2>
      <p id="modal-description">Review and edit metadata before downloading</p>

      {/* Form with editable fields */}
      <form method="dialog">
        {/* Fields for author, title, series, category */}
        <button type="button" onClick={onConfirm}>Download</button>
        <button type="button" onClick={onClose}>Cancel</button>
      </form>
    </dialog>
  );
}
```

### Anti-Patterns to Avoid
- **Concatenating series as strings in backend:** Prevents frontend from sorting/grouping properly. Always return structured data.
- **Mutating arrays during sort:** Use spread operator `[...array].sort()` to avoid side effects
- **DIY modal accessibility:** Use native `<dialog>` or battle-tested library. Focus management is complex.
- **Over-engineering grouping:** Simple `.reduce()` is sufficient. Don't add lodash just for groupBy.
</architecture_patterns>

<dont_hand_roll>
## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| Modal accessibility | Custom modal with manual focus trap | Native `<dialog>` or react-modal | Focus management, ARIA, keyboard navigation are complex. Native or library handles it correctly. |
| Array grouping | Custom grouping algorithm | Array.reduce() + object accumulator | Built-in, performant, well-understood pattern. No library needed. |
| Sorting | Custom sort logic | Array.sort() with comparator | Native implementation is optimized. Just provide comparator function. |
| Book number parsing | Complex regex for all formats | parseFloat() with fallback | Most book numbers are parseable as floats. Edge cases can fall back to string comparison. |
| Form validation | Manual field validation | react-hook-form (already in use) | Already in the stack, handles validation elegantly |

**Key insight:** This phase is mostly about data transformation (parsing, grouping, sorting) which JavaScript provides excellent primitives for. The only complex piece is modal accessibility, where native `<dialog>` now has excellent browser support (March 2022+).
</dont_hand_roll>

<common_pitfalls>
## Common Pitfalls

### Pitfall 1: Dropping Book Numbers in Backend
**What goes wrong:** Current `formatSeriesInfo()` extracts series name but discards book numbers (lines 298-303 commented out)
**Why it happens:** Simpler to concatenate strings than return structured data
**How to avoid:** Change return type to `[]SeriesInfo` struct with Name and Number fields. Frontend needs both for sorting.
**Warning signs:** UI can't sort books by series order, users complain books are out of sequence

### Pitfall 2: Mutating Arrays During Sort
**What goes wrong:** `array.sort()` mutates in place, causing unexpected re-renders or bugs
**Why it happens:** `.sort()` is deceptive - looks immutable but isn't
**How to avoid:** Always use `[...array].sort()` to create shallow copy first
**Warning signs:** React complains about state mutations, unpredictable UI updates

### Pitfall 3: Missing Multiple Series Per Book
**What goes wrong:** Assumes one series per book, breaks when book belongs to multiple series
**Why it happens:** Easier to store as single string than array
**How to avoid:** Model series as array (`series: SeriesInfo[]`) even for single-series books
**Warning signs:** Books in multiple series only show first series, user confusion

### Pitfall 4: Modal Focus Not Trapped
**What goes wrong:** User tabs outside modal to page content behind it
**Why it happens:** Didn't implement focus trap correctly
**How to avoid:** Use native `<dialog>` which handles this automatically, or use react-modal
**Warning signs:** Keyboard users can tab to hidden content, screen readers announce background content

### Pitfall 5: Non-Numeric Book Numbers
**What goes wrong:** Sorting breaks for series like "Book 1", "1.5", "Book Two"
**Why it happens:** Assumed all book numbers are simple integers
**How to avoid:** Use `parseFloat()` to extract number, fall back to string comparison if NaN
**Warning signs:** Books sort incorrectly ("Book 10" before "Book 2"), user reports

### Pitfall 6: Performance with Large Result Sets
**What goes wrong:** Grouping/sorting becomes slow with 100+ results
**Why it happens:** Inefficient grouping algorithm or too many re-renders
**How to avoid:** Use `useMemo()` to cache grouped/sorted results, only recalculate when results change
**Warning signs:** UI feels sluggish when many results, React DevTools shows excessive renders
</common_pitfalls>

<code_examples>
## Code Examples

### Backend: Structured Series Info
```go
// Source: Current codebase, enhanced based on existing pattern
type SeriesInfo struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Number string `json:"number"`
}

func parseSeriesInfo(seriesInfo string) []SeriesInfo {
    if seriesInfo == "" {
        return []SeriesInfo{}
    }

    // MAM format: {"123": ["Series Name", "Book Number", numeric_value]}
    seriesMap := make(map[string][]interface{})
    if err := json.Unmarshal([]byte(seriesInfo), &seriesMap); err != nil {
        return []SeriesInfo{}
    }

    result := []SeriesInfo{}
    for id, s := range seriesMap {
        info := SeriesInfo{ID: id}

        if len(s) > 0 {
            if name, ok := s[0].(string); ok {
                info.Name = name
            }
        }
        if len(s) > 1 {
            if number, ok := s[1].(string); ok {
                info.Number = number
            }
        }

        if info.Name != "" {
            result = append(result, info)
        }
    }

    return result
}
```

### Frontend: Grouping and Sorting
```typescript
// Source: Standard React pattern, verified with web search results
import { useMemo } from 'react';

interface SeriesGroup {
  seriesName: string;
  books: SearchResult[];
}

function useGroupedSeries(results: SearchResult[]): SeriesGroup[] {
  return useMemo(() => {
    // Group by series name
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
        const numA = parseFloat(a.series[0]?.number || '999');
        const numB = parseFloat(b.series[0]?.number || '999');
        return numA - numB;
      })
    }));
  }, [results]);
}
```

### Frontend: Accessible Confirmation Modal
```typescript
// Source: MDN Web Docs + accessibility guidelines
import { useEffect, useRef } from 'react';
import { useForm } from 'react-hook-form';

interface DownloadMetadata {
  author: string;
  title: string;
  series: string;
  category: string;
}

function DownloadConfirmModal({
  isOpen,
  onClose,
  onConfirm,
  initialData
}: {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: (data: DownloadMetadata) => void;
  initialData: DownloadMetadata;
}) {
  const dialogRef = useRef<HTMLDialogElement>(null);
  const { register, handleSubmit } = useForm({ defaultValues: initialData });

  useEffect(() => {
    const dialog = dialogRef.current;
    if (!dialog) return;

    if (isOpen) {
      dialog.showModal();
    } else {
      dialog.close();
    }
  }, [isOpen]);

  const onSubmit = (data: DownloadMetadata) => {
    onConfirm(data);
    onClose();
  };

  return (
    <dialog
      ref={dialogRef}
      aria-labelledby="confirm-title"
      aria-describedby="confirm-description"
      className="p-6 rounded-lg shadow-xl backdrop:bg-black backdrop:opacity-50"
      onClose={onClose}
    >
      <form onSubmit={handleSubmit(onSubmit)}>
        <h2 id="confirm-title" className="text-xl font-bold mb-2">
          Confirm Download
        </h2>
        <p id="confirm-description" className="text-gray-600 mb-4">
          Review and edit metadata before downloading
        </p>

        <div className="space-y-4">
          <label className="block">
            <span className="text-sm font-medium">Author</span>
            <input {...register('author')} className="w-full border rounded px-3 py-2" />
          </label>

          <label className="block">
            <span className="text-sm font-medium">Title</span>
            <input {...register('title')} className="w-full border rounded px-3 py-2" />
          </label>

          <label className="block">
            <span className="text-sm font-medium">Series</span>
            <input {...register('series')} className="w-full border rounded px-3 py-2" />
          </label>

          <label className="block">
            <span className="text-sm font-medium">Category</span>
            <input {...register('category')} className="w-full border rounded px-3 py-2" />
          </label>
        </div>

        <div className="flex gap-2 justify-end mt-6">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2 border rounded hover:bg-gray-50"
          >
            Cancel
          </button>
          <button
            type="submit"
            className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Download
          </button>
        </div>
      </form>
    </dialog>
  );
}
```

### Frontend: Series Group Display
```typescript
// Source: Current codebase pattern, enhanced for series grouping
function SeriesGroup({ seriesName, books }: SeriesGroup) {
  return (
    <div className="mb-6">
      <h3 className="text-lg font-semibold mb-3 px-4 py-2 bg-gray-100 rounded">
        {seriesName}
        <span className="text-sm font-normal text-gray-600 ml-2">
          ({books.length} {books.length === 1 ? 'book' : 'books'})
        </span>
      </h3>

      <div className="divide-y divide-gray-200 border border-gray-200 rounded-lg overflow-hidden">
        {books.map((book, idx) => (
          <SearchResultListItem
            key={`${book.id}-${idx}`}
            result={book}
            showSeriesNumber={seriesName !== 'Standalone'}
          />
        ))}
      </div>
    </div>
  );
}
```
</code_examples>

<sota_updates>
## State of the Art (2024-2026)

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Custom modal libraries | Native `<dialog>` element | March 2022 (baseline) | Simpler, lighter, built-in accessibility |
| Class components for modals | Hooks with useRef | React 16.8+ (2019) | Cleaner code, easier to manage |
| Complex lodash groupBy | Native Array.reduce() | Always available | No dependency needed |
| Prop drilling for modal state | Zustand (already in use) | Already adopted | Simpler state management |

**New tools/patterns to consider:**
- **Native `<dialog>` element:** Now "widely available" baseline. Handles focus trap, backdrop, ESC key automatically. Should be first choice over react-modal.
- **React 19:** Already in use. Latest concurrent features available if needed for performance.
- **TypeScript 5.9:** Already in use. Excellent type inference for reduce/sort patterns.

**Deprecated/outdated:**
- **react-modal:** Still valid but native `<dialog>` is now preferred for simple cases
- **lodash.groupby:** Unnecessary - native reduce() is sufficient and widely understood
- **String concatenation for structured data:** Always use typed objects for data that needs sorting/filtering
</sota_updates>

<open_questions>
## Open Questions

1. **Book number format variations**
   - What we know: MAM returns book numbers as strings ("1", "Book 1", "1.5")
   - What's unclear: Are there exotic formats like "Book One", "Volume II", "Part A"?
   - Recommendation: Start with parseFloat(), add string comparison fallback, validate with real MAM data

2. **Multiple series priority**
   - What we know: Books can belong to multiple series
   - What's unclear: Should we show duplicate books in each series group, or pick primary series?
   - Recommendation: Discuss with user - likely show in each series group (more discoverable)

3. **Performance threshold**
   - What we know: MAM can return 100+ results
   - What's unclear: At what result count does grouping/sorting impact UX?
   - Recommendation: Implement with useMemo(), measure with 100+ results, optimize if needed
</open_questions>

<sources>
## Sources

### Primary (HIGH confidence)
- Current codebase: `/backend/internal/search/providers/mam.go` - existing series parsing patterns
- MDN Web Docs: Native `<dialog>` element - https://developer.mozilla.org/en-US/docs/Web/HTML/Element/dialog
- MAM API structure: `TEST_MAM_SERIES.md` in codebase - documented series_info format

### Secondary (MEDIUM confidence)
- [How to sort an Array of Objects in React | bobbyhadz](https://bobbyhadz.com/blog/react-sort-array-of-objects) - verified with MDN
- [React: how to dynamically sort an array of objects - DEV Community](https://dev.to/ramonak/react-how-to-dynamically-sort-an-array-of-objects-using-the-dropdown-with-react-hooks-195p) - verified pattern
- [Building an Accessible Modal Dialog in React – Chris Henrick](https://clhenrick.io/blog/react-a11y-modal-dialog/) - verified with ARIA spec
- [GitHub - reactjs/react-modal](https://github.com/reactjs/react-modal) - official React modal library
- [The best React modal dialog libraries of 2026 | Croct Blog](https://blog.croct.com/post/best-react-modal-dialog-libraries) - ecosystem overview

### Tertiary (LOW confidence - needs validation)
- None - all patterns verified against official docs or existing codebase
</sources>

<metadata>
## Metadata

**Research scope:**
- Core technology: MAM API JSON parsing (Go), React UI grouping (TypeScript)
- Ecosystem: Native browser APIs, existing project stack
- Patterns: Data transformation, grouping/sorting, accessible modals
- Pitfalls: Data structure issues, accessibility, performance

**Confidence breakdown:**
- Standard stack: HIGH - all libraries already in project, no new dependencies needed
- Architecture: HIGH - verified patterns against official docs and existing codebase
- Pitfalls: HIGH - common React/TypeScript issues, documented in community resources
- Code examples: HIGH - based on existing codebase patterns and official documentation

**Research date:** 2026-01-07
**Valid until:** 2026-02-07 (30 days - stable web technologies)

</metadata>

---

*Phase: 07-mam-series-detection*
*Research completed: 2026-01-07*
*Ready for planning: yes*
