---
id: TH-0078
title: "I-0007: how should the review limit stop breaking close-outs?"
anchor:
  path: wip/kanban/stories/S-0243-review-column-exceeds-its-wip-limit-while-acceptance-is-batched.md
  item: S-0243
status: resolved
participants: [agent-S-0243, alex]
created: 2026-10-03T02:44:57Z
updated: 2026-10-03T03:26:35Z
---

# TH-0078 I-0007: how should the review limit stop breaking close-outs?

On wip/kanban/stories/S-0243-review-column-exceeds-its-wip-limit-while-acceptance-is-batched.md.

## Entries

### 2026-10-03T02:44:57Z agent-S-0243
From S-0243. I-0007's cause still occurs, and the issue leaves the fix to you. Two things combine today:

1. Pulling ignores the review limit. `flai serve`, `inbox` (can_pull) and `wait_for_work` look only at the in-progress limit (serve/agents.go:330), so while acceptance is batched, review grows without bound (5 today, 3 in progress).
2. `flai check --strict` fails on `board.wip-limit`. Once review is over its limit, every other story's `scripts/close-out.sh` stops on a warning its agent cannot clear; only acceptance clears it. workflow.md says a limit breach "does not block", but `--strict` makes it block.

Raising the limit (the issue's first option) only moves the threshold, so it doesn't fix the cause.

Recommended (A), both halves:
- The pull respects the review limit: while review is at or over its limit, serve starts no new story, and inbox and wait_for_work offer none. The ready story is held with the reason "review is full (5 of 5): accept or send back a story". Stories already in progress still finish.
- The review column over its limit stays a warning in `flai check` and on the board, but `--strict` no longer fails on it: it is your queue, not the story's. The ready and in-progress limits keep failing `--strict` as they do now.
- Tests: (a) check --strict passes with review over its limit; (b) serve, inbox and wait_for_work hold a ready story while review is full.

B: only the pull hold. A story already in progress that finishes into a full review still stops its close-out, as today.
C: only the strict change. Review keeps growing while acceptance is batched, and the board shows it over the limit.

Which do you want? I'll update workflow.md and the docs to match.

### 2026-10-03T02:57:35Z alex
A

### 2026-10-03T03:26:35Z agent-S-0243
Resolved: Option A, built in S-0243 and recorded in ADR-0073.
