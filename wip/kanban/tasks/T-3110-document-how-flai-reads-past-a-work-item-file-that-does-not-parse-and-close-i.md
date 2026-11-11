---
id: T-3110
type: task
nature: remediation
title: Document how flai reads past a work item file that does not parse, and close I-0344
status: backlog
parent: S-0716
owner: alex
created: 2026-11-11T08:56:59Z
updated: 2026-11-11T08:56:59Z
transitions: []
stream: S-0716
tags: [docs, issues]
touches: [design/system/work-hierarchy.md, docs/users/flai.md, design/issues/I-0344-one-malformed-work-item-file-in-the-main-checkout-s-wip-makes-every-flai-read.md, design/issues/summary.md]
after: [T-3108, T-3109]
---
# T-3110 Document how flai reads past a work item file that does not parse, and close I-0344

## Work

In `design/system/work-hierarchy.md` § Front matter schema, next to § Fields an older flai does not know, say what happens to a work item file that does not parse. Every read leaves it out and goes on, and `flai check` reports it as one finding on its path, by the rule T-3108 named. In `docs/users/flai.md` § Check the repository, add that rule to the list of what the check covers. If the story's agent judges the loader's new contract an architectural decision, record it in an ADR with `flai adr new` and link it from both places.

Then close the issue with `flai issue close I-0344 --reason`. Name the fix: `Repo.List` skips the file, `flai check` reports it, and the tests T-3107, T-3108, and T-3109 added.

This task waits for T-3108, whose rule name the docs give, and T-3109, whose tests the close reason cites.

## Done when

- Both documents describe the behaviour T-3107 and T-3108 built, and their `updated` dates are bumped.
- I-0344 is closed with a reason naming the fix, and `design/issues/summary.md` no longer lists it.
- `flai check --strict` and the markdown lint pass on the changed files.

## Notes
