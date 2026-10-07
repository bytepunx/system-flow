---
id: T-1184
type: task
nature: improvement
title: The user guide and the dashboard's design say times are shown in local time, and a sweep finds no UTC display left
status: backlog
parent: S-0329
owner: alex
created: 2026-10-07T19:50:05Z
updated: 2026-10-07T19:50:15Z
transitions: []
stream: S-0329
tags: [dashboard, docs]
touches: [docs/users/flaiover.md, design/system/flaiover-dashboard.md]
after: [T-1181, T-1182, T-1183]
---
# T-1184 The user guide and the dashboard's design say times are shown in local time, and a sweep finds no UTC display left

## Work

- `docs/users/flaiover.md`: say once, where a reader looks first (the header or the Board section), that the dashboard shows every time in the browser's time zone while flai records them in UTC. Under Charts, keep that days and weeks are UTC buckets, since `flai stats` counts them so, and say an hour is labelled in local time: the line `Days and hours are UTC.` under Spend over time, the Agent waiting bar's `Monday to Sunday, UTC`, and the strategic charts' `per UTC day` are read again against what T-1183 built. The schedules' `cron expression in UTC` rows stay: they say how the expression is read.
- `design/system/flaiover-dashboard.md`: record the rule under Internal structure, or the section that holds the shared helpers: times arrive from flai as `YYYY-MM-DDTHH:MM:SSZ` and are shown in the browser's zone through `flaiover/src/lib/localtime.ts`, and no component formats a time of its own. Bump `updated`.
- Sweep `flaiover/src` for a time still shown in UTC: `UTC` in rendered text, `useUTC`, `.slice(0, 10)` or `.slice(0, 16)` on a time, and a raw `.at`, `started`, `ended`, `created`, or `published` in markup. Fix what the view tasks missed, adding the file to this task's touches, or note in the narrative why it stays.

It waits for the three view tasks because it describes what they built and sweeps after them.

## Done when

- `docs/users/flaiover.md` and `design/system/flaiover-dashboard.md` say times are shown in local time and how the chart buckets relate to it.
- The sweep finds no time shown in UTC in `flaiover/src`, and both of the story's criteria are verified and ticked.
- The markdown lint passes on both documents.

## Notes

- Out of scope, as the plan thread says: the raw front matter a document page shows in its collapsed Front matter panel, and times written inside a markdown document's body, such as a narrative's log headings. Both are the file's own text.
- S-0328, in progress when this was planned, also changes `docs/users/flaiover.md`. Its claim holds this story in ready until it leaves in-progress; once it is accepted, `flai stream sync` brings its edits in before this task writes.
