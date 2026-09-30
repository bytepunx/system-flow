---
id: T-0606
type: task
nature: improvement
title: The item page offers New story or New epic, starting the form with the same type
status: done
parent: S-0171
owner: alex
created: 2026-09-30T01:12:22Z
updated: 2026-09-30T01:16:07Z
transitions:
  - to: ready
    at: 2026-09-30T01:12:25Z
    by: agent-S-0171
  - to: in-progress
    at: 2026-09-30T01:12:25Z
    by: agent-S-0171
  - to: done
    at: 2026-09-30T01:16:07Z
    by: agent-S-0171
stream: S-0171
tags: []
touches: [flaiover/src]
usage:
  source: log
  seconds: 222
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 49
      output: 10917
      cache_read: 2320537
      cache_write: 59544
      cost: 1.159
---
# T-0606 The item page offers New story or New epic, starting the form with the same type

## Work

- On a story's or an epic's page, when the dashboard may write, show a `New story` or `New epic` link at the top, beside the title, that opens `/new?type=<type>`.
- A story's link also carries its epic (`&parent=E-nnnn`) so the new story starts under the same epic; `/new` passes it to `NewItemForm` as where the parent select starts.
- No link on a task's page: tasks are the agent's to write, and the dashboard makes only epics and stories.
- Behavior tests: the link's text and href for a story and an epic, none for a task or when read-only; the form starts with the type and parent it was given.

## Done when

- The link shows at the top of story and epic pages with the viewed type in its text.
- Following it opens the form set to that type, and a story's to its epic.
- `npm run check`, `npm run lint`, and `npm test` pass in `flaiover/`.

## Notes

- The link shows on archived, done, and cancelled items too: making another of the same type does not depend on the one viewed being open.
- A story made from a story whose epic has since closed starts with no epic, since the form offers open epics only.
- The link carries a query string after `resolve('/new')`, which `svelte/no-navigation-without-resolve` cannot follow; it is disabled around the link with that reason, as the board and the layout do.
- Verified with `scripts/flaiover-test.sh`: prettier, eslint, svelte-check, and 668 vitest tests pass. The Playwright tests were not run.
