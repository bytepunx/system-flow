---
id: S-025
type: story
nature: feature
title: Tooling support for conventions
status: backlog
parent: E-005
owner: alex
created: 2026-09-15T18:26:32Z
updated: 2026-09-15T18:26:32Z
transitions: []
tags: [conventions, cli]
---

# S-025 Tooling support for conventions

## Goal
flai and the dashboard understand the conventions folder so it is validated, discoverable, and easy to load.

## Acceptance criteria
- [ ] system-flow.yaml layout gains `conventions`; manifest.Load requires it; template.yaml layout and the prototype manifest include it
- [ ] flai check validates conventions front matter and that README.md lists every convention file exactly once
- [ ] `flai prime` prints the conventions in read order (paths by default, `--cat` for full content) so an agent or a hook can load them in one call
- [ ] flai upgrade design (S-020) notes that convention files merge above their marker like CLAUDE.md
- [ ] design/system/flaiover-dashboard.md adds conventions to the documentation explorer and search scope

## Tasks

## Notes
