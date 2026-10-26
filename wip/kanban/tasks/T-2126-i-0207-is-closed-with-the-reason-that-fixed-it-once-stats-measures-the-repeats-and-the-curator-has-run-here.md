---
id: T-2126
type: task
nature: remediation
title: I-0207 is closed with the reason that fixed it, once stats measures the repeats and the curator has run here
status: backlog
parent: S-0504
owner: alex
created: 2026-10-26T04:10:42Z
updated: 2026-10-26T04:10:42Z
transitions: []
stream: S-0504
tags: [issues, digest]
touches: [design/issues/I-0207-the-curator-has-never-run-on-system-flow-so-no-thread-is-in-any-entry-s-asked-again-and-the-digest-s-repeat-measure-reads-0-whatever-was-asked.md, design/issues/summary.md]
after: [T-2122, T-2124, T-2125]
---
# T-2126 I-0207 is closed with the reason that fixed it, once stats measures the repeats and the curator has run here

## Work

With the curator's first run recorded (T-2122) and `flai stats` telling an unmeasured repeat count from a measured one (T-2123, T-2125, T-2124), verify the cause no longer occurs: run `scripts/flai.sh stats --json --since 7d` in the story's worktree and read `digest.curator` and the current week's `curator_runs` and `asked_again`, which is present once a run has reached the week. Record the figures in the narrative's `## Decisions`.

Then close the issue from the worktree, so the close is committed on the story branch:

```bash
scripts/flai.sh issue close I-0207 --reason "<what fixed it>"
```

The reason names both halves: the operator turned `curate` on for system-flow on <date> and the first run ended at <time>, and `flai stats` now reports `digest.weeks[].curator_runs` and `digest.curator` and leaves `asked_again` out of a week no run has reached (the ADR's ID). `flai issue close` writes the reason under the issue's Remediation and regenerates `design/issues/summary.md`; commit both with the task.

Tick the story's criteria with `flai criteria tick S-0504 <n>` for each you verified. This waits on T-2122, T-2124, and T-2125 because the issue is closed only once both halves of the remedy are in place and documented; T-2123 is covered through T-2125.

## Done when

- `flai stats --json` on this project reports `digest.curator` with at least one run and the current week's `asked_again`, and the figures are in the narrative's `## Decisions`.
- I-0207 reads `status: closed` with a reason under Remediation that names the action's enabling, the first run, and the ADR, and `design/issues/summary.md` no longer lists it as open.
- Both acceptance criteria of S-0504 are ticked through flai.

## Notes

Drafted by planner-S-0504. The re-run of S-0435's measurement from 2026-11-09, which the experiment's recommendation names, is not this story's: it needs two whole weeks after the curator's first run and is proposed in the plan's thread as a story of its own.
