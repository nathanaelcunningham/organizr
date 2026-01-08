# Phase 7: MAM Series Detection - Context

**Gathered:** 2026-01-07
**Status:** Ready for research

<vision>
## How This Should Work

When I search MAM and get results, books that belong to a series should be visually grouped together by series name, with the books ordered by their series number within each group. This makes it easy to scan and see complete series or find specific books in order.

Before downloading, I want to review and confirm (or edit if needed) the metadata that will be sent to the backend: author, title, series information, and category. This ensures everything is correct before the download starts.

MAM results can include books that belong to multiple series, so the parsing and display needs to handle that complexity gracefully.

</vision>

<essential>
## What Must Be Nailed

- **Accurate series parsing** - Extracting series name and book number correctly from MAM search results
- **Proper grouping and ordering** - Visual organization showing series grouped together with books in numerical order
- **Handling multiple series** - Books that belong to multiple series need to be parsed and displayed correctly
- **Pre-download confirmation** - Ability to review and edit author, title, series, and category before submitting to backend

All of these are equally important for this phase to work properly.

</essential>

<boundaries>
## What's Out of Scope

- **Series metadata enrichment** - Not looking up additional series info from external sources; only use what MAM provides
- **Filtering by series** - Not adding search filters or controls to show/hide specific series; just display what's in the results

</boundaries>

<specifics>
## Specific Ideas

- **Clean visual separation** - Series groups should be clearly distinct and easy to scan, making it obvious which books belong together
- **Confirmation flow** - Open to best practices for the pre-download review/edit experience (modal, inline, or form approach)

</specifics>

<notes>
## Additional Context

The core user experience is about making series relationships obvious in the search results, so users can quickly understand what they're looking at and find the books they want. The confirmation step ensures accuracy before committing to a download.

Multiple series per book is a known complexity that must be addressed in the parsing logic.

</notes>

---

*Phase: 07-mam-series-detection*
*Context gathered: 2026-01-07*
