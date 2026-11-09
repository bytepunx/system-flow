---
id: T-2803
type: task
nature: improvement
title: Measure 13 board.get calls at once before and after S-0602 and S-0606 on the operator's host, and record them under a new section of server-performance.md
status: backlog
parent: S-0636
owner: alex
created: 2026-11-09T22:07:12Z
updated: 2026-11-09T22:07:21Z
transitions: []
stream: S-0636
tags: [cli]
touches: [design/system/server-performance.md]
---
# T-2803 Measure 13 board.get calls at once before and after S-0602 and S-0606 on the operator's host, and record them under a new section of server-performance.md

## Work

Measure as S-0604's T-2681 did, so the figures compare with § Cause 14 fixed in `design/system/server-performance.md`.

- Build two flai trees in the same hour, under `/tmp`. Before is `main` at `5e1b4a1c8^`, the commit before S-0602's acceptance: it has none of S-0602, S-0604, or S-0606. After is `main` at the story's base, with all three.
- Write an in-process Go program under `/tmp`, against each tree. It calls `board.get` through `hostapi.MethodsFor` with a request timer, for the projects in `~/.flai/serve/projects.json`. It runs a watcher over `mainRefs` for each root, as flai serve does. The after build also calls `release.Watch` and `release.RefsMoved`, as S-0604's build did.
- Each run makes a warm-up call and one sequential call per project. Then it makes five rounds of 13 calls released together. Run the two builds in turn, twice each.
- Give the phases too: `release.pending`, `merge.queue`, `repo.list`, `undone.read`, and `board.view`, for all projects and for this repository.
- Record the method, the host's one-minute load, the project count, how many of this repository's stories were in review, and the before and after. Put them in a new section after § Cause 14 fixed, named for S-0636. Link it from that section's last paragraph, which names S-0636.

It waits for no task: it is the first layer.

## Done when

- [ ] `design/system/server-performance.md` has a section for S-0636 with 13 `board.get` calls at once, before and after, per project and per phase, and how they were measured.
- [ ] The section says whether every call of every project answered under 300 ms. Where one did not, it names the project, the phase that took the time, and why.

## Notes
