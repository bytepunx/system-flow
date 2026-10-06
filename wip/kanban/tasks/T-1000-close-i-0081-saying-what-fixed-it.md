---
id: T-1000
type: task
nature: remediation
title: Close I-0081 saying what fixed it
status: backlog
parent: S-0284
owner: alex
created: 2026-10-06T06:24:25Z
updated: 2026-10-06T06:24:25Z
transitions: []
stream: S-0284
tags: [issues]
touches: [design/issues/I-0081-flai-s-permission-prompt-waits-for-an-answer-from-the-story-s-owner-and-the-operator-s-thread-replies-carry-another-name-so-their-allow-is-never-seen-and-the-write-times-out.md, design/issues/summary.md]
after: [T-0998, T-0999]
---
# T-1000 Close I-0081 saying what fixed it

## Work

Once T-0998's fix and T-0999's documents are in, close the issue from the story's worktree: `flai issue close I-0081 --reason "..."`, naming the ADR from T-0997 and the change to `awaitAnswer`, that an answer may come from the project's owner as well as the story's owner, and the test that reproduces S-0218's instance. Check that `design/issues/summary.md` shows it closed.

## Done when

- I-0081 is closed with a reason naming the ADR, the fix, and its test
- `design/issues/summary.md` lists I-0081 as closed
- `flai check --strict` is clean

## Notes
