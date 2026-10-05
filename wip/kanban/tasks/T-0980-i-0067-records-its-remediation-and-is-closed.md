---
id: T-0980
type: task
nature: remediation
title: I-0067 records its remediation and is closed
status: backlog
parent: S-0254
owner: alex
created: 2026-10-05T05:49:48Z
updated: 2026-10-05T05:49:48Z
transitions: []
stream: S-0254
tags: [issues]
touches: [design/issues/I-0067-flai-touches-with-paths-replaces-an-item-s-touches-and-its-help-does-not-say-so.md, design/issues/summary.md]
after: [T-0978]
---
# T-0980 I-0067 records its remediation and is closed

## Work

Write I-0067's `## Remediation`: `flai touches` gained `--add` and `--remove`, its help says that paths given alone replace the list, and `flai stream sync` hints with `--add` (S-0254, T-0976, T-0978). Then close it from the story's worktree, since `flai issue` writes under the checkout it runs in: `flai issue close I-0067 --reason "S-0254: flai touches --add and --remove, and its help says paths alone replace the list"`. That keeps `design/issues/summary.md` current.

Waits for T-0978, so that the issue is closed only when the fix is built and documented.

## Done when

- I-0067's `## Remediation` names the fix and S-0254
- `flai issue list` no longer shows I-0067 as open, and `design/issues/summary.md` follows
- `flai check --strict` scoped to the story is clean

## Notes
