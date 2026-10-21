---
id: T-1959
type: task
nature: improvement
title: settingspages.ts maps the five Configuration pages to their sections, and SettingsPanel renders only the sections of the page it is told
status: done
parent: S-0459
owner: alex
created: 2026-10-21T07:03:04Z
updated: 2026-10-21T09:04:18Z
transitions:
  - to: ready
    at: 2026-10-21T08:51:42Z
    by: system-flow
  - to: in-progress
    at: 2026-10-21T08:51:43Z
    by: system-flow
  - to: done
    at: 2026-10-21T09:04:18Z
    by: agent-S-0459
stream: S-0459
tags: [dashboard]
touches: [flaiover/src/lib/settingspages.ts, flaiover/src/lib/settingspages.test.ts, flaiover/src/lib/components/SettingsPanel.svelte, flaiover/src/lib/components/SettingsPanel.svelte.test.ts]
usage:
  source: log
  seconds: 755
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 96
      output: 478
      cache_read: 5135011
      cache_write: 167145
      cost: 1.9767
---
# T-1959 settingspages.ts maps the five Configuration pages to their sections, and SettingsPanel renders only the sections of the page it is told

## Work

- Write `flaiover/src/lib/settingspages.ts`. It holds the five pages, each with a slug (`features`, `agents`, `file-system`, `projects`, `harnesses`), a title, and its section keys, as TH-0553 question 4 settles them. It also maps each old `/settings` anchor to its page.
- Give `SettingsPanel.svelte` a `page` prop. Wrap each of its ten sections in a test of whether the page holds it. Keep the load, the settings gate, the read-only notice, and the save state as one script.
- Test in `settingspages.test.ts` that every section belongs to exactly one page and every slug is unique.
- Run each existing case in `SettingsPanel.svelte.test.ts` with the `page` its section is on. Add one case per page that checks only its sections render.
- Waits for nothing in this story, which waits for S-0458.

## Done when

Both test files pass, and `SettingsPanel` with each `page` renders exactly that page's sections.

## Notes
