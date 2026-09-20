---
id: T-0298
type: task
nature: feature
title: An item's page has an edit mode for stories and epics, and criteria can be ticked in place
status: done
parent: S-0085
owner: alex
created: 2026-09-20T15:21:42Z
updated: 2026-09-20T15:36:26Z
transitions:
  - to: ready
    at: 2026-09-20T15:36:25Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T15:36:26Z
    by: system-flow
  - to: done
    at: 2026-09-20T15:36:26Z
    by: system-flow
stream: S-0085
tags: []
---
# T-0298 An item's page has an edit mode for stories and epics, and criteria can be ticked in place

## Work
Fields as inputs (title, nature, tags, touches for stories, parent for stories), the body as Markdown with a preview, what stays flai's shown and not editable. Save shows a refusal's findings and keeps the text, shows a conflict with the choice to load or overwrite. Acceptance criteria tick in place from the page. Tasks are not editable here and the page says so.

## Done when
- Component tests for the form, a refusal, a conflict, and a ticked criterion
- The flaiover tests and build pass

## Notes
