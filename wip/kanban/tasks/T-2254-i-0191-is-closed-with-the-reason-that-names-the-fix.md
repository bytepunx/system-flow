---
id: T-2254
type: task
nature: remediation
title: I-0191 is closed with the reason that names the fix
status: backlog
parent: S-0525
owner: alex
created: 2026-10-28T14:53:03Z
updated: 2026-10-28T14:53:03Z
transitions: []
stream: S-0525
tags: [flai, issues]
touches: [design/issues/I-0191-a-work-item-write-in-a-project-with-no-gitignore-makes-flai-cache-items-lock-and-then-refuses-itself-on-layout-gitignore.md, design/issues/summary.md]
after: [T-2251]
---
# T-2254 I-0191 is closed with the reason that names the fix

## Work

From the story's worktree, on its branch, once T-2251's fix and test are committed, close the issue in one call that commits what it writes:

```sh
flai issue close I-0191 --commit --reason "S-0525: a checked creation scopes its check to the new item's story as a checked edit does, so layout.gitignore on the repository root warns instead of refusing the first write; the items lock stays under .flai-cache"
```

It rewrites the issue's file and `design/issues/summary.md`, commits both on `story/S-0525`, and widens the story's touches. Then tick the story's second criterion with `flai criteria tick S-0525 2`, and the first once T-2251's tests pass. It waits for T-2251 because the reason names a fix that must exist.

## Done when

- `design/issues/I-0191-…md` has status `closed` with the reason above, and `design/issues/summary.md` no longer lists it among the open issues.
- The close is committed on the story branch, and `flai issue list` no longer shows I-0191.

## Notes

`flai issue close` takes no `--story`; the story resolves from the worktree's branch.
