---
id: T-1528
type: task
nature: research
title: Set out where and in what format a project keeps its design system, and how the dashboard previews it live, with the options and a recommendation
status: done
parent: S-0380
owner: alex
created: 2026-10-13T17:46:01Z
updated: 2026-10-13T18:07:38Z
transitions:
  - to: ready
    at: 2026-10-13T18:00:56Z
    by: agent-S-0380
  - to: in-progress
    at: 2026-10-13T18:00:56Z
    by: agent-S-0380
  - to: done
    at: 2026-10-13T18:07:38Z
    by: agent-S-0380
stream: S-0380
tags: [dashboard, cli]
touches: [design/system/design-system-tool.md]
after: [T-1527]
usage:
  source: log
  seconds: 402
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 38
      output: 13417
      cache_read: 1904299
      cache_write: 66194
      cost: 1.0496
---
# T-1528 Set out where and in what format a project keeps its design system, and how the dashboard previews it live, with the options and a recommendation

## Work

Add to the finding the options for the format and the place of a project's design system: tokens in a markdown document's front matter under `design/system`, a folder of its own under `design`, a W3C design tokens JSON file, or a Tailwind v4 stylesheet as the source, each weighed against the documentation standard, the host methods, `flai check`, and what an agent reads. Then the options for previewing tokens live without building the project: an iframe or a shadow root carrying CSS custom properties over sample components compiled with the dashboard, Tailwind's browser build, or a build on the host; each weighed for isolation from the dashboard's theme, speed, and what the container may run (ADR-0031). End each with a recommendation. It waits for T-1527, whose section it builds on in the same file.

## Done when

- The finding has a section on format and place and a section on the live preview, each with its options compared in a table and a recommendation.

## Notes

- Layer 2 of 5.
