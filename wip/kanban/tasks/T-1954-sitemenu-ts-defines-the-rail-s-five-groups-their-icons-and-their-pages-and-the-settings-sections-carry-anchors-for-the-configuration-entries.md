---
id: T-1954
type: task
nature: feature
title: sitemenu.ts defines the rail's five groups, their icons, and their pages, and the settings sections carry anchors for the Configuration entries
status: done
parent: S-0458
owner: alex
created: 2026-10-21T07:01:42Z
updated: 2026-10-21T08:03:17Z
transitions:
  - to: ready
    at: 2026-10-21T07:57:05Z
    by: agent-S-0458
  - to: in-progress
    at: 2026-10-21T07:57:06Z
    by: agent-S-0458
  - to: done
    at: 2026-10-21T08:03:17Z
    by: agent-S-0458
stream: S-0458
tags: [dashboard]
touches: [flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts, flaiover/src/lib/components/SiteMenu.svelte, flaiover/src/lib/components/SiteMenu.svelte.test.ts]
after: [T-1953]
usage:
  source: log
  seconds: 371
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 90
      output: 429
      cache_read: 3386056
      cache_write: 103265
      cost: 1.2967
---
# T-1954 sitemenu.ts defines the rail's five groups, their icons, and their pages, and the settings sections carry anchors for the Configuration entries

## Work

- Replace `SITE_MENU` in `flaiover/src/lib/sitemenu.ts` with the five groups in the epic's order: Flow, Agents, Design, Documentation, and Configuration. Each group has a name and an icon name. Each page has a label, an icon name, a route through `resolve()`, and an optional badge.
- Use the routes in the story's second criterion. Mark Workflow `soon`, with no route.
- Keep `locate`, the longest-path match, and `groupBadges`. Move `/items`, `/review`, and `/new` to Flow and `/edit` to Documentation.
- Give the five Configuration sections' headings in `SettingsPanel.svelte` ids for the entries' anchors: `features`, `agents`, `file-system`, `projects`, and `harnesses`, on the first section of each as TH-0553 question 4 groups them.
- Update `sitemenu.test.ts`:
  - every route `locate` maps, including the documents under `/docs`;
  - Workflow has no link;
  - every page has an icon.
- Icon names are strings here, so this waits for nothing.

## Done when

`sitemenu.test.ts` passes on the five groups, and each Configuration anchor exists on `/settings`.

## Notes
