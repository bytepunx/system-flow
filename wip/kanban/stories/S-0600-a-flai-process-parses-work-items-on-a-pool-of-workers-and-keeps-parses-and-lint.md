---
id: S-0600
type: story
nature: improvement
title: A flai process parses work items on a pool of workers and keeps parses and lint results under .flai-cache, so a cold list takes under 0.3 s, not 1.1 to 1.5 s
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:39:35Z
updated: 2026-11-07T22:39:35Z
transitions: []
tags: [cli]
touches: [flai/internal/workitem/store.go, flai/internal/workitem/transition.go, flai/internal/mdlint]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0600 A flai process parses work items on a pool of workers and keeps parses and lint results under .flai-cache, so a cold list takes under 0.3 s, not 1.1 to 1.5 s

## Goal

Cause 10 in [server-performance.md](../../../design/system/server-performance.md#where-the-time-goes): S-0156's item cache lives in flai serve's memory, so every `flai` a write starts, every `flai check`, and every `flai move` parses all 3,233 items one file after another (`flai/internal/workitem/store.go`), 1.1 to 1.5 s, and lints all of `wip` again, 0.98 s. `flai move` reads the archive only so that `Follow` sees an epic's archived stories (`flai/internal/workitem/transition.go`); a move takes 1 to 2.5 s, 80% of it the list.

## Acceptance criteria

- [ ] `store.list` parses on a bounded pool of workers, and the markdown lint of a folder lints its files the same way.
- [ ] A flai process keeps parsed items, and the lint findings of each file, under `.flai-cache`, keyed by path, size, and modification time and by the lint configuration's hash, and a file whose key differs is read again. A missing, damaged, or foreign cache is a miss, never an error, as S-0531's kept commits are.
- [ ] A cold `flai move` and a cold `flai check` in a scratch copy of this repository each take at least 1 s less than before, and server-performance.md records the before and after.

## Tasks

## Notes

- Measured: one worker 1.28 s, 8 workers 0.65 s, 16 workers 0.49 s on the loaded host; the lint of `wip` 0.98 s to 0.25 s on 8 workers.
- The cache is rebuildable from the files and changes nothing in `wip/` or `design/`. It is written only into a `.flai-cache` that exists, as D-0402 says.
