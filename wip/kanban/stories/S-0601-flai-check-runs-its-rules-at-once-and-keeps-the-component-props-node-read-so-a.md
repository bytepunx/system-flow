---
id: S-0601
type: story
nature: improvement
title: flai check runs its rules at once and keeps the component props node read, so a whole check takes about 1.5 s, not 3.3 to 8 s
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:39:49Z
updated: 2026-11-07T22:39:49Z
transitions: []
tags: [cli]
touches: [flai/internal/check/check.go, flai/internal/components/props.go]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0601 flai check runs its rules at once and keeps the component props node read, so a whole check takes about 1.5 s, not 3.3 to 8 s

## Goal

Cause 11 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): `check.Run` runs its rules one after another (`flai/internal/check/check.go`), and `componentExamples` starts `node` on every run (`flai/internal/components/props.go`), 0.6 to 4.2 s. A whole check takes 3.3 s cold and 1.8 s warm without `node`, and 7.4 to 8 s cold and 2.6 s warm with it. `flai check`, `flai verify`, the close-out's check step, and `flai task done` pay it.

## Acceptance criteria

- [ ] `check.Run` runs its independent rules at once, gathering their findings under one lock, and gives the same findings in the same order as before, which a test over a fixture project with a finding of every rule shows.
- [ ] The props `node` reads are kept under `.flai-cache`, keyed by each component file's path, size, and modification time and by the sub-project's svelte version, and `node` starts only for a component whose key changed.
- [ ] A whole `flai check --strict` on this repository, with nothing changed since the last, takes under 2 s, and server-performance.md records the before and after.

## Tasks

## Notes

- The rules only read, so running them at once needs no other change. Run the tests under `-race`.
