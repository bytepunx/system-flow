---
id: T-0459
type: task
nature: remediation
title: Thread entries render as markdown on every page that shows threads
status: done
parent: S-0126
owner: alex
created: 2026-09-26T07:31:26Z
updated: 2026-09-26T07:33:30Z
transitions:
  - to: ready
    at: 2026-09-26T07:31:29Z
    by: agent-S-0126
  - to: in-progress
    at: 2026-09-26T07:31:29Z
    by: agent-S-0126
  - to: done
    at: 2026-09-26T07:33:30Z
    by: agent-S-0126
stream: S-0126
tags: [dashboard]
touches: [flaiover/src/lib]
---
# T-0459 Thread entries render as markdown on every page that shows threads

## Work

- Render each entry of `flaiover/src/lib/components/Threads.svelte` through `render` from `$lib/markdown`, in a `prose` block, instead of a pre-wrapped paragraph. Threads.svelte is the one component that shows entries on an item page, the review page, the document viewer, and the document editor.
- Resolve relative links in an entry from the repository root, since entry authors write repository paths.
- A component test that mounts Threads with an entry holding a list, emphasis, inline code, and a link, and finds the elements rather than the raw markdown.

## Done when

- The test passes, and fails against the old pre-wrapped paragraph.
- `npm run lint`, `npm run check`, and the unit suite pass in `flaiover/`.

## Notes
