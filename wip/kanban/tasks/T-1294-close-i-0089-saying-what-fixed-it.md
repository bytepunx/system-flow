---
id: T-1294
type: task
nature: remediation
title: Close I-0089 saying what fixed it
status: done
parent: S-0315
owner: alex
created: 2026-10-07T23:44:02Z
updated: 2026-10-07T23:57:29Z
transitions:
  - to: ready
    at: 2026-10-07T23:57:26Z
    by: agent-S-0315
  - to: in-progress
    at: 2026-10-07T23:57:26Z
    by: agent-S-0315
  - to: done
    at: 2026-10-07T23:57:29Z
    by: agent-S-0315
stream: S-0315
tags: [cli]
touches: [design/issues/I-0089-design-issues-summary-md-is-generated-and-committed-so-a-story-branch-that-records-an-issue-conflicts-with-any-issue-recorded-on-main-meanwhile.md, design/issues/summary.md]
after: [T-1293]
usage:
  source: log
  seconds: 3
  estimated: true
  models: []
---
# T-1294 Close I-0089 saying what fixed it

## Work

Close I-0089 in the story's worktree with `flai issue close I-0089 --reason "<reason>"`, which writes the reason under its Remediation section and regenerates `design/issues/summary.md` without it. The reason names what fixed it and what shows it: S-0278's regeneration of the summary when a sync's or an acceptance's rebase stops on it alone ([ADR-0098](../../../design/adrs/0098-flai-stream-sync-and-flai-accept-regenerate-design-issues-summary-md-when-a.md)), and the test T-1293 added, which reproduces an issue recorded and one closed on main while a story branch records its own.

It waits for T-1293: the reason cites the test, and the issue is closed only once the test shows the cause no longer occurs. If T-1293 found a cause ADR-0098 does not cover, the reason names the fix task instead.

## Done when

- I-0089 has `status: closed` and its Remediation section carries the reason.
- `design/issues/summary.md` no longer lists I-0089.
- `flai check --strict` is clean for the story.

## Notes

Drafted by the planner for S-0315.
