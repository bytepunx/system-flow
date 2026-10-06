---
id: T-1020
type: task
nature: remediation
title: A flai serve test reproduces the idle board behind a story in review and shows the orchestrator's acceptance under accept_reviews clears it
status: done
parent: S-0296
owner: alex
created: 2026-10-06T11:50:06Z
updated: 2026-10-06T20:15:41Z
transitions:
  - to: ready
    at: 2026-10-06T20:11:19Z
    by: agent-S-0296
  - to: in-progress
    at: 2026-10-06T20:11:19Z
    by: agent-S-0296
  - to: done
    at: 2026-10-06T20:15:41Z
    by: agent-S-0296
stream: S-0296
tags: [flai]
touches: [flai/internal/serve/review_wait_test.go]
usage:
  source: log
  seconds: 262
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 45
      output: 14544
      cache_read: 2237652
      cache_write: 85785
      cost: 1.2876
---
# T-1020 A flai serve test reproduces the idle board behind a story in review and shows the orchestrator's acceptance under accept_reviews clears it

## Work

Write `flai/internal/serve/review_wait_test.go`, beside `hold_test.go` and `orchestrate_test.go`, reproducing I-0088 in a scratch project. Since ADR-0096 (S-0295) a story in review holds nothing by overlap, and `TestTheLauncherStartsAStoryOverlappingOnlyOneInReview` shows it. What still waits on its acceptance is a ready story that names it in `after:`, and every ready story while review is at its limit. So the scratch project has one story in review with every criterion ticked and its branch verified, review at a limit of one, a ready story naming it in `after:`, and a ready story with no relation to it.

- With `orchestration.permissions.accept_reviews` off, the orchestrator's acceptance is refused, the launcher starts no story agent, and both ready stories wait: this is the cause the issue describes.
- With `accept_reviews` on, a stand-in orchestrator accepts the story in review through the path S-0221 built for `flai accept --by orchestrator --verified --evidence` (the permission, then `preview.AcceptWith` with the orchestrator's options), and the launcher then starts both waiting stories' agents without a person.

Use the stand-in harness the existing serve tests use. A test that needs real git skips under `-short`. If the launcher notices the acceptance only on its one-minute look, so that the test has to wait for it, make it look when a story leaves review, in `orchestrate.go`. The task waits for nothing in this story; S-0221, which the story waits for, provides the orchestrator's acceptance.

## Done when

- The test fails with `accept_reviews` off, the board idle behind the story in review, and passes with it on, the waiting stories started.
- `scripts/flai-test.sh` passes, the race detector included.

## Notes
