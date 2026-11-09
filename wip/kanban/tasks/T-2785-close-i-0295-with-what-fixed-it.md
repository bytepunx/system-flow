---
id: T-2785
type: task
nature: improvement
title: Close I-0295 with what fixed it
status: backlog
parent: S-0628
owner: alex
created: 2026-11-09T21:59:27Z
updated: 2026-11-09T21:59:27Z
transitions: []
stream: S-0628
tags: [issues]
touches: [design/issues/I-0295-the-close-out-job-and-the-manifest-s-test-commands-run-the-story-worktree-s.md, design/issues/summary.md]
after: [T-2772, T-2777, T-2782, T-2784]
---
# T-2785 Close I-0295 with what fixed it

## Work

Run `flai issue close I-0295 --reason "<reason>"`. The reason says that a container story's close-out job now runs its verify and tiers in the story's container (T-2777), and that the operator's checks do too (T-2782). It says the host's commit and sync run flai's own code with the main checkout's hooks. It names what still runs on the host by design, as ADR-0272 records it: the merge queue's batch, an acceptance's `fix_command`, and `flai test` or `flai verify` run by hand. It also says that the Impact's `manifest.go:1570` is `KlusterErrors`, and that the tiers are `TestTiers` there. It names the test that reproduces the cause. flai regenerates `design/issues/summary.md`.

Waits for every other task of the story, since the reason names what they did.

## Done when

- I-0295 is closed with that reason, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` passes.

## Notes
