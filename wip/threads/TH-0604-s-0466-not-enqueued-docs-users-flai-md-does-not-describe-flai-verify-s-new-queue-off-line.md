---
id: TH-0604
title: "S-0466 not enqueued: docs/users/flai.md does not describe flai verify's new queue_off line"
anchor:
  path: wip/kanban/stories/S-0466-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md
  item: S-0466
status: resolved
participants: [orchestrator, alex]
created: 2026-10-23T05:12:24Z
updated: 2026-10-23T17:43:05Z
---

# TH-0604 S-0466 not enqueued: docs/users/flai.md does not describe flai verify's new queue_off line

On wip/kanban/stories/S-0466-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md.

## Entries

### 2026-10-23T05:12:24Z orchestrator
**Recommended:** send S-0466 back to in-progress for one docs edit, then it goes to the queue.

- In `docs/users/flai.md` § Tiers left to the merge queue (around line 1624), add a sentence. When a tier marked `queue` runs because the queue cannot batch, `flai verify` names the condition that is off and the configuration file it read.
- Add `queue_off` to the `--json` field list (around line 1593).
- Add `docs/users/flai.md` to the story's touches.

## What I checked

`flai verify S-0466` passed at the branch head, 59ae9c36. Integration and smoke are queued, which shows that TH-0601's configuration works. The verifier matched both criteria to changed files within the touches:

- **Criterion 1:** `queueOff` in `flai/cmd/verify.go`, `Report.QueueOff` (`queue_off`) in `flai/internal/verify/story.go`, and three cases in `flai/cmd/verify_test.go`. The dev docs are in `docs/contributors/index.md` and `design/system/devex.md`.
- **Criterion 2:** I-0148 is closed with its reason.

## The blocker

The project addition to `design/conventions/documentation.md` says `docs/users/flai.md` is updated in the same story that changes a command. This story changes `flai verify`'s text and JSON output, and `docs/users/flai.md` is neither in the diff nor in the touches. Only the generated reference changed.

### 2026-10-23T17:43:05Z alex
Resolved: S-0466 was accepted
