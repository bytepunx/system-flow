---
id: T-2054
type: task
nature: remediation
title: I-0148 is closed with the reason naming flai close-out and its detached run
status: cancelled
parent: S-0466
owner: alex
created: 2026-10-22T21:01:34Z
updated: 2026-10-23T03:32:56Z
transitions:
  - to: cancelled
    at: 2026-10-23T03:32:56Z
    by: agent-S-0466
stream: S-0466
tags: [cli]
touches: [design/issues/I-0148-a-close-out-whose-run-main-outruns-takes-two-rounds-of-the-integration-and-smoke-tiers-about-24-minutes-past-the-bash-tool-s-10-minute-limit-so-the-agent-s-call-is-killed-mid-round.md, design/issues/summary.md]
after: [T-2050]
---
# T-2054 I-0148 is closed with the reason naming flai close-out and its detached run

## Work

In the story's worktree, once the command, the scripts, and the docs are in, close the issue with `flai issue close I-0148 --reason "<reason>"`, the reason saying what fixed it: the close-out runs detached in `flai close-out`, the agent's call follows it for its wait and is called again, so a run that main outruns finishes its second round whatever the Bash tool does to the call, and naming the ADR T-2030 wrote. The command rewrites the issue file and regenerates `design/issues/summary.md`; commit both. Tick the story's criteria you verified with `flai criteria tick`, and put the `design/issues/summary.md` table in the review report, as `communication.md` asks.

Waits for T-2050, the last task that changes behaviour or its description: the issue is closed only when the fix is whole.

## Done when

- `design/issues/I-0148-…md` has `status: closed` and a `## Remediation` entry with the reason, and `design/issues/summary.md` no longer lists it among the open issues.
- `flai check --strict` is clean, and the story's second criterion is ticked.

## Notes

The fifth and last layer.
- 2026-10-23T03:32:56Z: moved to cancelled: TH-0594 settled on option 3: S-0466 makes the close-out see the merge queue; the detached close-out is left to a backlog story from I-0148
