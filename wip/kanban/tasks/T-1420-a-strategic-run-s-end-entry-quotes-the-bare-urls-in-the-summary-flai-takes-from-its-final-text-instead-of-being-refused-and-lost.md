---
id: T-1420
type: task
nature: remediation
title: A strategic run's end entry quotes the bare URLs in the summary flai takes from its final text, instead of being refused and lost
status: backlog
parent: S-0346
owner: alex
created: 2026-10-08T08:54:03Z
updated: 2026-10-08T08:54:03Z
transitions: []
stream: S-0346
tags: [flai, serve]
touches: [flai/internal/serve/activity.go, flai/internal/serve/activity_test.go, design/system/agent-narrative.md]
after: [T-1419]
---
# T-1420 A strategic run's end entry quotes the bare URLs in the summary flai takes from its final text, instead of being refused and lost

## Work

- `logRunEndSaying` in `flai/internal/serve/activity.go` takes the summary from the first line of the run's final text, which no agent checks against the lint. When that line holds a bare URL, `updateActivity` refuses the entry through `LintGuard`, and `flai serve` only warns. The run's seconds and cost are then lost from the activity document. This is the one path left by which I-0118's bare URL reaches `wip/agents/<kind>.md` from a flai that has S-0324's rule.
- Pass the summary that `say` returns through `mdlint.QuoteBareURLs` before it is appended, for the planner, the orchestrator, and the analyzer alike. Leave `activity_log` as it is: an agent's own summary is still refused with the rule and the line, as `tooling.md` says.
- Reproduce it in `flai/internal/serve/activity_test.go`: a run whose final text starts with a line holding a bare `www.` literal is logged, with the literal in a code span, and the document lints clean. Without the fix the entry is refused.
- Say in `design/system/agent-narrative.md` § Strategic agents' activity documents that flai quotes a bare URL in the summary it takes from a run's last result.
- Waits for T-1419, whose `QuoteBareURLs` it calls.

## Done when

- The reproducing test fails without the change and passes with it.
- `flai test flai/internal/serve` passes.
- `design/system/agent-narrative.md` says what the run-end summary does with a bare URL, with `updated` bumped.

## Notes
