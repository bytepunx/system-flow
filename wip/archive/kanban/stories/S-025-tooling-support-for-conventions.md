---
id: S-025
type: story
nature: feature
title: Tooling support for conventions
status: done
parent: E-005
owner: alex
created: 2026-09-15T18:26:32Z
updated: 2026-09-16T23:14:35Z
transitions:
  - to: ready
    at: 2026-09-16T23:03:59Z
    by: agent
  - to: in-progress
    at: 2026-09-16T23:04:00Z
    by: agent
  - to: review
    at: 2026-09-16T23:07:56Z
    by: agent
  - to: done
    at: 2026-09-16T23:14:35Z
    by: alex
tags: [conventions, cli]
---

# S-025 Tooling support for conventions

## Goal
flai and the dashboard understand the conventions folder so it is validated, discoverable, and easy to load.

## Acceptance criteria
- [x] flai resolves the folder as `<layout.design>/conventions`; the template prototype ships it and `flai new` renders it
- [x] flai check validates conventions front matter and that README.md lists every convention file exactly once
- [x] `flai prime` prints the conventions in read order (paths by default, `--cat` for full content) so an agent or a hook can load them in one call
- [x] flai upgrade design (S-020) notes that convention files merge above their marker like CLAUDE.md
- [x] design/system/flaiover-dashboard.md adds conventions to the documentation explorer and search scope

## Tasks
- T-051 conventions loader and check rules
- T-052 flai prime command
- T-053 Fixtures, docs, and design notes for conventions tooling

## Notes
- Check rules caught two empty untracked files in the template conventions (logging, telemetry); scaffolded as drafts for the operator.
