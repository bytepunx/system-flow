---
id: T-0512
type: task
nature: feature
title: CLAUDE.md and session-start.md prime with --story when there is a story
status: done
parent: S-0138
owner: alex
created: 2026-09-29T01:08:52Z
updated: 2026-09-29T01:13:30Z
transitions:
  - to: ready
    at: 2026-09-29T01:08:56Z
    by: agent-S-0138
  - to: in-progress
    at: 2026-09-29T01:12:30Z
    by: agent-S-0138
  - to: done
    at: 2026-09-29T01:13:30Z
    by: agent-S-0138
stream: S-0138
tags: [template]
touches: [CLAUDE.md, template/root/CLAUDE.md.tmpl, template/root/design/conventions/session-start.md, design/conventions/session-start.md, design/system/conventions.md, template/template.yaml, template/CHANGELOG.md]
---
# T-0512 CLAUDE.md and session-start.md prime with --story when there is a story

## Work

Say in `CLAUDE.md`, `template/root/CLAUDE.md.tmpl`, and the baseline `session-start.md` (template, then copied here above the marker) to prime with `flai prime --story <id>` when there is a story and `flai prime --cat` otherwise, and to read what the catalog lists when it is needed. Update `design/system/conventions.md` "Priming" and drop the not-yet-built note. Bump the template's version and record it in its changelog.

## Done when

- The four files and the design say the same; the template version and changelog record it; `make smoke` passes.

## Notes
