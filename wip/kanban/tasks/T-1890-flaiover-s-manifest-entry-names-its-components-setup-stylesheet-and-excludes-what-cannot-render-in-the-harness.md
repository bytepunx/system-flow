---
id: T-1890
type: task
nature: feature
title: flaiover's manifest entry names its components' setup stylesheet and excludes what cannot render in the harness
status: backlog
parent: S-0444
owner: alex
created: 2026-10-16T21:39:16Z
updated: 2026-10-16T21:39:59Z
transitions: []
stream: S-0444
tags: [dashboard]
touches: [system-flow.yaml]
---
# T-1890 flaiover's manifest entry names its components' setup stylesheet and excludes what cannot render in the harness

## Work

Start the dev server on flaiover and open each component in the harness. Find which fail for want of the app's runtime: its server data, `$app` state with no page, or the channel. Set flaiover's `components` with `flai components config flaiover --setup src/routes/layout.css --exclude <glob>...`, never by hand. Write each exclusion's reason in the story's `## Decisions` for T-1893 to carry into the design.

Waits for nothing in this story; it needs S-0437 to S-0443 accepted.

## Done when

`flai components list --project flaiover` lists everything but the exclusions, the harness renders each listed component or shows only a missing-prop error, and `flai check --strict` passes.

## Notes

`CLAUDE.md` says flai owns the manifest but for `name`, `description`, and `dashboard`, so the write goes through `flai components config` (S-0437).
