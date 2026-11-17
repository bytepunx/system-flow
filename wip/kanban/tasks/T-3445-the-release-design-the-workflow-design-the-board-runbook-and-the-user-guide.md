---
id: T-3445
type: task
nature: feature
title: The release design, the workflow design, the board runbook, and the user guide describe the suites a publish runs
status: backlog
parent: S-0811
owner: alex
created: 2026-11-16T23:48:23Z
updated: 2026-11-16T23:48:23Z
transitions: []
stream: S-0811
tags: [flai, docs]
touches: [design/system/pushing-from-the-board.md, design/system/workflow.md, docs/operators/runbooks/change-the-board.md, docs/users/flai.md]
after: [T-3441, T-3443]
---
# T-3445 The release design, the workflow design, the board runbook, and the user guide describe the suites a publish runs

## Work

Describe what T-3441 and T-3443 built, linking ADR-0312 §6 and repeating no paragraph of it:

- `design/system/pushing-from-the-board.md`: in the opening current-state paragraph, that with `verify_at.publish` set `flai release --pending` runs those suites at the main branch's head once something is pending, before it applies a version or tags, and on a failure publishes nothing and exits with T-3441's code naming the suite, the tier, and the log.
- `design/system/workflow.md`, in the workflow file's section on `verify_at`: what `publish` runs, when, how the tiers are selected, and that the baseline names none.
- `docs/operators/runbooks/change-the-board.md`: a section on running suites before a publish: the key to add, what a failure looks like, where its log is, and what to do next.
- `docs/users/flai.md`: the publishing section gains the run and the new exit, and the publish formula paragraph names the publish step's new `stop` exit.

It waits for T-3441, whose exit and output it describes, and for T-3443, whose `stop` it names.

## Done when

- The four documents say the above, and each has its `updated` date moved.
- The markdown lint and `flai check --strict` pass on them.

## Notes
