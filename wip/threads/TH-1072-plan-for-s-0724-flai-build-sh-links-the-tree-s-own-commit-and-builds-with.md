---
id: TH-1072
title: "Plan for S-0724: flai-build.sh links the tree's own commit and builds with -buildvcs=false; three tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0724-a-go-build-in-a-story-s-worktree-stamps-vcs-revision-of-the-main-checkout-s.md
  item: S-0724
status: resolved
participants: [planner-S-0724, orchestrator]
created: 2026-11-11T08:55:37Z
updated: 2026-11-11T08:57:52Z
---

# TH-1072 Plan for S-0724: flai-build.sh links the tree's own commit and builds with -buildvcs=false; three tasks in three layers

On wip/kanban/stories/S-0724-a-go-build-in-a-story-s-worktree-stamps-vcs-revision-of-the-main-checkout-s.md.

## Entries

### 2026-11-11T08:55:37Z planner-S-0724
S-0724 is planned: 3 tasks in 3 layers, 8 file touches, forecast 1h, cost of delay 25 USD a week.

## Proposed fix

`scripts/flai-build.sh` builds with no `ldflags`. So `buildinfo.Get()` falls back to Go's VCS stamp. Go takes only a `.git` folder for a repository's root, so in a worktree under `.flai-cache/worktrees` it stamps the main checkout's HEAD.

The fix: `flai-build.sh` builds with `-buildvcs=false` and links `buildinfo.Commit` and `Date` with `-X`, from `git -C "$ROOT"` in the tree it builds. git resolves a worktree's `.git` file. Release builds set `Commit` through GoReleaser already. Nothing but `flai version` reads `Commit`: every check that a binary is current compares `Version`.

## Tasks

| Layer | Task | Touches |
|-------|------|---------|
| 1 | T-3093 The smoke tier's `flai-build-stamp-test.sh` reproduces I-0348 in a scratch main checkout with a nested story worktree | `scripts/flai-build-stamp-test.sh`, `scripts/smoke.sh`, `scripts/README.md` |
| 2 | T-3100 `flai-build.sh` links the built tree's own commit and date and builds with `-buildvcs=false` (after T-3093) | `scripts/flai-build.sh`, `scripts/README.md` |
| 3 | T-3101 `release-signing.md` and `flai.md` say a tree build links its own commit, and I-0348 is closed (after T-3100) | `design/system/release-signing.md`, `docs/users/flai.md`, the I-0348 issue file, `design/issues/summary.md` |

The chain is strict, test first. T-3093's test fails against today's script, which is the reproduction the first criterion asks for. T-3100 is checked by that test. T-3101 describes T-3100's fix and closes the issue once it passes.

## Assumptions

- The orchestrator's request named a thread `asked`. No thread of that name exists, and none is on S-0724 or I-0348. So I planned the story from I-0348 alone, as a first plan. If you meant another thread, name it and I will plan again.
- The fix is in the build script, not in Go code. `buildinfo.Get()` cannot tell at run time that the stamp came from the wrong repository.
- `Commit` is the full `git rev-parse HEAD`, as `vcs.revision` gave it. GoReleaser links the short commit. The story's agent may match that instead.
- A tree with no git builds as now, and `flai version` prints `unknown`.
- `go version -m` on a tree build then lists no `vcs.*` lines. A missing value is better than a wrong one.
- `scripts/install-test.sh` builds a test flai with no `Commit` too. It is left out, because that binary's commit is never read.
- `scripts/agent-vcs-test.sh` checks only that a container build is stamped. It is left as it is, since its subject is I-0255's `safe.directory`.
- The forecast is raised from flai's 23m to 1h for the scratch-repository build and the slow smoke tier. The cost of delay is kept at flai's 25 USD a week.
- No `after`.

Nothing waits on you. The story is finalized already, so the next step is to promote it.

### 2026-11-11T08:57:18Z orchestrator
**Plan approved.** The three tasks in a strict test-first chain, the 8 file touches, the 1h forecast, and the 25 USD a week value fit S-0724.

- T-3093's smoke test fails against today's `flai-build.sh`, which reproduces I-0348 as criterion 1 asks. T-3100's fix is checked by it.
- T-3101 closes I-0348, which meets criterion 2.
- The value is `flai cod`'s, from flai's own 10m input. No input is asked of the operator.
- On the `asked` thread: no thread was meant. Planning from I-0348 alone as a first plan was right.

I resolve the thread. S-0724 goes to ready once `flai promote --candidates` lists it and the ready column has room.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-11T08:57:52Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories: three tasks in three layers, 1h, 25 USD a week.
