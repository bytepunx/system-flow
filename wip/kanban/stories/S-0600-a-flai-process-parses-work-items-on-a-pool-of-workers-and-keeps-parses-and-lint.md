---
id: S-0600
type: story
nature: improvement
title: A flai process parses work items on a pool of workers and keeps parses and lint results under .flai-cache, so a cold list takes under 0.3 s, not 1.1 to 1.5 s
status: backlog
parent: E-0033
owner: alex
created: 2026-11-07T22:39:35Z
updated: 2026-11-07T23:49:19Z
transitions: []
tags: [cli]
touches: [flai/internal/workitem/store.go, flai/internal/workitem/transition.go, flai/internal/mdlint, flai/internal/workitem/store_test.go, flai/internal/workitem/kept.go, flai/internal/workitem/kept_test.go, flai/internal/mdlint/config.go, flai/internal/mdlint/folder.go, flai/internal/mdlint/folder_test.go, flai/internal/mdlint/kept.go, flai/internal/mdlint/kept_test.go, flai/internal/check/check.go, design/system/server-performance.md, design/system/flai-cli.md, docs/operators/runbooks/backup.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
finalized:
  by: alex
  at: 2026-11-07T23:49:17Z
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

### Planning

Touches:

| Touch | From | Why |
|-------|------|-----|
| `flai/internal/workitem/store.go` | declared | `store.list` loops over files at lines 61 to 96; the pool and the kept parses hook in here |
| `flai/internal/workitem/transition.go` | declared | `TransitionAll` lists the archive at line 70 and lists again at line 130 for the index; the second list can come from the first |
| `flai/internal/mdlint` | declared | Kept as declared. The folder lint and the kept findings are new files in it, which the tasks name; flai narrows it to them in the claim |
| `flai/internal/workitem/store_test.go` | co-change (26%) | The store's tests count parses per file and cover the pool |
| `flai/internal/workitem/kept.go`, `kept_test.go` | layout | New: the items kept under `.flai-cache`, modelled on `flai/internal/storygit/committed.go` |
| `flai/internal/mdlint/config.go` | layout | `Load` keeps no bytes and no hash; the key needs one |
| `flai/internal/mdlint/folder.go`, `folder_test.go` | layout | New: the folder lint on a pool; today the loop is in `check.go` |
| `flai/internal/mdlint/kept.go`, `kept_test.go` | layout | New: the lint findings kept per file |
| `flai/internal/check/check.go` | layout | `markdown()` at lines 1898 to 1928 walks `wip` and lints one file at a time |
| `design/system/server-performance.md` | design | Criterion 3 records the before and after |
| `design/system/flai-cli.md` | co-change (13%) | `## Internal structure` says what flai keeps and where |
| `docs/operators/runbooks/backup.md` | layout | Its table lists each `.flai-cache` file and whether it is worth keeping |

`docs/users/flai.md` (26% co-change) is left out: no command, flag, or setting changes. The worker count is bounded by the CPUs, with no setting.

No `after`: the story reads nothing another open story writes. S-0598 and S-0601 change `flai/internal/check/check.go` too; a trial merge judges that.

Forecast and cost of delay: below, once written.
