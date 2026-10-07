---
id: T-1140
type: task
nature: remediation
title: I-0100 is closed with flai issue close, saying that acceptance retries a held index lock and finishes a failed commit on a second run
status: backlog
parent: S-0307
owner: alex
created: 2026-10-07T01:12:57Z
updated: 2026-10-07T01:12:57Z
transitions: []
stream: S-0307
tags: [flai]
touches: [design/issues/I-0100-the-acceptance-s-commit-step-fails-when-another-process-holds-git-s-index-lock-and-the-acceptance-cannot-then-be-finished-by-flai.md, design/issues/summary.md]
after: [T-1138]
---
# T-1140 I-0100 is closed with flai issue close, saying that acceptance retries a held index lock and finishes a failed commit on a second run

## Work

Close I-0100 from the story's worktree with `flai issue close I-0100 --reason "<what fixed it>"`, which writes the issue and regenerates `design/issues/summary.md`. The reason names S-0307 and the fix: the acceptance's commit waits out an index lock another process holds, and an acceptance done and archived whose commit failed is finished by running `flai accept` again. It waits for T-1138, since the reason states what that task delivered, and runs beside T-1139, with which it shares no path.

## Done when

- I-0100's front matter says closed, and its Remediation says what fixed it and names S-0307.
- `design/issues/summary.md` no longer lists I-0100.
- `flai check --strict` reports nothing on either file.

## Notes
