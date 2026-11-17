---
id: T-3556
type: task
nature: improvement
title: The design and the user docs describe what a rebase and resolve job owes a moved main and how it names its rounds, and I-0398 is closed
status: backlog
parent: S-0824
owner: alex
created: 2026-11-17T08:34:10Z
updated: 2026-11-17T08:34:10Z
transitions: []
stream: S-0824
tags: [docs, flai, flaiover]
touches: [design/system/flai-cli.md, design/system/flaiover-dashboard.md, docs/users/flai.md, docs/users/flaiover.md, design/issues/I-0398-a-rebase-and-resolve-job-runs-the-whole-close-out-again-each-time-main-moves.md, design/issues/summary.md]
after: [T-3554, T-3555]
---
# T-3556 The design and the user docs describe what a rebase and resolve job owes a moved main and how it names its rounds, and I-0398 is closed

## Work

Write down what T-3554 and T-3555 built, and close the issue. It waits for both, since the docs describe their behaviour and the issue closes once the fix is in.

- `design/system/flai-cli.md` § The close-out job: in the Rebase and resolve and state bullets, say that a resolve job's last check passes over a main that moved outside the story's change set, linking the new ADR. Say what the round fields of the job's record and `resolve.status` hold.
- `docs/users/flai.md` § The close-out job and § Resuming across a clean rebase: say when a resolve job ends in the round it verified and when it runs another, and that the acceptance's `verify_behind_main` judges what main brings in later.
- `design/system/flaiover-dashboard.md` and `docs/users/flaiover.md` § Reviewing a story: say the review page shows the round and why.
- Close I-0398 with `flai issue close I-0398 --reason "<what fixed it>"`, naming the ADR, the flag, and the tests T-3554 added. This regenerates `design/issues/summary.md`. Tick criterion 2 once it is closed, and criterion 1 once T-3554's tests pass.

## Done when

- The four documents describe the behaviour, and the markdown lint passes on them.
- I-0398 is closed with a reason saying what fixed it, and `design/issues/summary.md` lists it closed.
- `flai check --strict` is clean for the story.

## Notes
