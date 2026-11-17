---
id: T-3552
type: task
nature: improvement
title: flai verify --sync-only, when asked, passes a branch main outran while verify_behind_main would raise no blocker at the head it verified
status: backlog
parent: S-0824
owner: alex
created: 2026-11-17T08:33:31Z
updated: 2026-11-17T08:33:31Z
transitions: []
stream: S-0824
tags: [flai, verify]
touches: [flai/internal/preview/accept.go, flai/internal/preview/accept_test.go, flai/internal/verify/story.go, flai/internal/verify/story_test.go, flai/cmd/verify.go, flai/cmd/verify_test.go, docs/users/flai-reference.md]
after: [T-3551]
---
# T-3552 flai verify --sync-only, when asked, passes a branch main outran while verify_behind_main would raise no blocker at the head it verified

## Work

Build the last check T-3551's ADR decides. It waits for T-3551, which names the flag and the condition.

Today `verifyBehindMain` in `flai/internal/preview/accept.go` judges the acceptance's gate:

- `mainPathsSinceVerify` lists the paths main brought in since the report's merge base;
- `overlappingMainPaths` keeps those that overlap the change set, through `verify.OverlappingPaths`;
- `verify.StaleTiers` names the tiers they send back.

`preview` imports `verify`, so `verify.SyncOnly` cannot call it. Move the judgement into package `verify`, or behind a function `flai/cmd/verify.go` calls after `SyncOnly`. The acceptance and the last check then share one rule, and `verifyBehindMainGate` keeps its words and its behaviour.

With the ADR's flag, `flai verify <story> --sync-only` judges a branch that lacks main's commits by that rule against the story's last passing report. It passes when no tier is sent back. Its last line then says main moved, by how many paths, and that none overlaps the change set or sends a tier back. Without the flag it stops at sync as today, and so does it when a tier is sent back, a stop that names the tiers and the first overlapping paths for the round line. A report, or a main, that cannot be read stops it, so the check never opens unread.

Tests:

- in `flai/internal/verify/story_test.go`, or `flai/cmd/verify_test.go` where the flag is read: a main that moved outside the change set passes with the flag and stops without it, and a main that moved a path the change set holds stops with the flag and names the tier;
- in `flai/internal/preview/accept_test.go`: `TestAnAcceptanceWaitsForAVerifyOverMainsNewPaths` still passes over the moved code.

Update `flai verify`'s help and `docs/users/flai-reference.md` § flai verify for the flag.

## Done when

- `flai verify --sync-only` with the ADR's flag passes over a main whose new paths send back no tier by ADR-0303's rule, and stops naming the tiers when they do.
- The acceptance's `verify_behind_main` and the last check share one implementation, and the acceptance's tests pass unchanged.
- The new tests pass with `flai test` on the changed paths.

## Notes
