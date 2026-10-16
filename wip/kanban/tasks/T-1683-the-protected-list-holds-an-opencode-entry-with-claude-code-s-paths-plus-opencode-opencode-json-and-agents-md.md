---
id: T-1683
type: task
nature: feature
title: The protected list holds an opencode entry with Claude Code's paths plus .opencode/, opencode.json, and AGENTS.md
status: backlog
parent: S-0406
owner: alex
created: 2026-10-16T05:14:03Z
updated: 2026-10-16T05:14:03Z
transitions: []
stream: S-0406
tags: [cli]
touches: [flai/internal/protected/protected.go, flai/internal/protected/protected_test.go]
---
# T-1683 The protected list holds an opencode entry with Claude Code's paths plus .opencode/, opencode.json, and AGENTS.md

## Work

- Add `opencode` to `lists` in `flai/internal/protected/protected.go`: the folders and files of the `claude-code` entry, plus `.opencode/`, `opencode.json`, and `AGENTS.md`, as ADR-0130 names them. Replace the comment that says `.opencode/` joins later.
- Test `For("opencode")` and `Path` on each new path, and that `claude-code`'s list is unchanged.
- Waits for nothing.

## Done when

`protected.For("opencode")` returns the list, its tests pass, and the `claude-code` list is as it was.

## Notes
