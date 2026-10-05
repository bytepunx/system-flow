---
id: T-0983
type: task
nature: remediation
title: Close I-0056 saying what fixed it
status: backlog
parent: S-0265
owner: alex
created: 2026-10-05T05:50:11Z
updated: 2026-10-05T05:50:11Z
transitions: []
stream: S-0265
tags: [issues]
touches: [design/issues/I-0056-flai-s-wip-markdown-lint-does-not-flag-a-bare-email-address-so-a-thread-entry-with-one-reached-main-and-failed-a-story-s-close-out.md, design/issues/summary.md]
after: [T-0981, T-0982]
---
# T-0983 Close I-0056 saying what fixed it

## Work

In the story's worktree, run `flai issue close I-0056 --reason` with a reason naming the fix: flai's MD034 reports GFM's extended email autolink as markdownlint-cli2 does, tested by the email fixture and `TestBareEmailOfI0056` (S-0265, T-0981).

Waits for T-0981, whose fix and test the reason names, and for T-0982, so that the issue closes only when the story's work is in.

## Done when

- I-0056 is closed with that reason, and `design/issues/summary.md` lists it closed.
- `flai check --strict` is clean.

## Notes
