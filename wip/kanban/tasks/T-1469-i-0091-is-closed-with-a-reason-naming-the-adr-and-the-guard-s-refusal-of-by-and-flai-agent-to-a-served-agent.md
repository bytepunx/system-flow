---
id: T-1469
type: task
nature: remediation
title: I-0091 is closed with a reason naming the ADR and the guard's refusal of --by and FLAI_AGENT to a served agent
status: backlog
parent: S-0369
owner: alex
created: 2026-10-12T17:06:09Z
updated: 2026-10-12T17:06:09Z
transitions: []
stream: S-0369
tags: [issues]
touches: [design/issues/I-0091-a-thread-entry-s-author-is-whatever-by-says-so-an-agent-with-a-shell-could-answer-a-permission-thread-as-the-owner.md, design/issues/summary.md]
after: [T-1468]
---
# T-1469 I-0091 is closed with a reason naming the ADR and the guard's refusal of --by and FLAI_AGENT to a served agent

## Work

Close I-0091 with `flai issue close I-0091 --reason "<reason>"`, run in the story's worktree. The reason names the ADR from T-1465, the guard's refusal from T-1467, and the tests that reproduce the cause. It also names the residuals the ADR leaves open, such as a shell redirection into a thread file.

It waits for T-1468, the last of the work the reason describes, and so for every other task.

## Done when

- I-0091's status is closed and its reason says what fixed it.
- `design/issues/summary.md` no longer lists it, and `flai check --strict` is clean.

## Notes

This ticks criterion 2. Criterion 1 is met by T-1467's tests.
