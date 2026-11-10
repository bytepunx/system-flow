---
id: T-2989
type: task
nature: remediation
title: Close I-0273 with the fix that takes the watcher's first snapshot before the first looks
status: backlog
parent: S-0690
owner: alex
created: 2026-11-10T10:29:54Z
updated: 2026-11-10T10:29:54Z
transitions: []
stream: S-0690
tags: [issues]
touches: [design/issues/I-0273-flai-serve-takes-a-project-s-watcher-s-first-snapshot-after-its-first-looks-so.md, design/issues/summary.md]
after: [T-2984]
---
# T-2989 Close I-0273 with the fix that takes the watcher's first snapshot before the first looks

## Work

It waits for T-2984, whose fix and test the reason names.

Run `flai issue close I-0273 --reason "<reason>"` in the story's worktree. The reason names the method T-2981 added, that flai serve now calls it before the item snapshot and the first looks (T-2984), and the serve test that reproduces the issue. `flai issue close` writes the issue file and regenerates `design/issues/summary.md`.

## Done when

- I-0273 is closed, with a reason that names the fix and its test.
- `design/issues/summary.md` no longer lists I-0273 as open.
- `flai check --strict` reports nothing on the issue.

## Notes
