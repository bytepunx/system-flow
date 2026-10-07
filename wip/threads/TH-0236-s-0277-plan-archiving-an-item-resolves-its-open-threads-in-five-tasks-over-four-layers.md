---
id: TH-0236
title: "S-0277 plan: archiving an item resolves its open threads, in five tasks over four layers"
anchor:
  path: wip/kanban/stories/S-0277-flai-check-finds-threads-archived-outside-the-story-at-close-out.md
  item: S-0277
status: resolved
participants: [planner-S-0277, alex]
created: 2026-10-07T01:21:35Z
updated: 2026-10-07T01:23:42Z
---

# TH-0236 S-0277 plan: archiving an item resolves its open threads, in five tasks over four layers

On wip/kanban/stories/S-0277-flai-check-finds-threads-archived-outside-the-story-at-close-out.md.

## Entries

### 2026-10-07T01:21:35Z planner-S-0277
**Plan.** Archiving an item, by `flai accept`, `flai move <id> done`, or `flai archive`, resolves every thread still `open` or `answered` on it. Each thread gets an entry saying why, in the same commit. That stops I-0073 at its source.

All 27 instances are threads left on a story that was then accepted. TH-0112, TH-0127, TH-0194, TH-0203, and TH-0232 each warned until someone ran `flai thread resolve` by hand. Today, only the orchestrator's acceptance looks at threads (ADR-0093).

Tasks and layers:

| Layer | Task | Waits for | Touches |
|-------|------|-----------|---------|
| 1 | T-1142 `threads.OnItems` and `ResolveOnItems`; the `threads.archived` message names `flai thread resolve` | none | `flai/internal/threads/archived*.go`, `flai/internal/check/check*.go` |
| 2 | T-1144 `flai accept` resolves them in the acceptance commit; the dry run lists them; the test reproduces I-0073 | T-1142 | `flai/cmd/accept.go`, `flai/internal/preview/accept.go`, `flai/cmd/accept_threads_test.go` |
| 2 | T-1146 `flai archive` does the same | T-1142 | `flai/cmd/archive.go`, `flai/cmd/archive_test.go` |
| 3 | T-1147 ADR, `flai-cli.md`, `workflow.md`, `docs/users/flai.md`, regenerated reference | T-1144, T-1146 | docs, `design/adrs` |
| 4 | T-1148 resolve the threads left over (TH-0203, TH-0232 today), `flai check --strict` clean, close I-0073 | T-1147 | `design/issues/I-0073-*`, `summary.md` |

T-1144 and T-1146 share no file and run together.

Assumptions:

- An acceptance by you resolves a thread that is still `open`, not only an `answered` one. Accepting the story is your word that its questions are settled. The dry run names those threads before you accept. The alternative is to refuse your acceptance while a thread is open, as the orchestrator's is refused. I'd rather not, because it adds a step to every board acceptance.
- The orchestrator's refusal while a thread is open stays as it is.
- No thread archive folder: threads stay in `wip/threads/`, resolved.
- The story is held while S-0286 is in progress, through `flai/cmd/accept.go` and `flai/internal/preview/accept.go`. It also overlaps S-0307 on those files.
- The cost of delay waits on TH-0233.

Proposals: none to split, merge, or drop.

### 2026-10-07T01:23:42Z alex
Resolved.
