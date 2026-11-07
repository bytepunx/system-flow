---
id: T-2593
type: task
nature: remediation
title: Close I-0303 with the reason naming the guardSettings fix and its test
status: backlog
parent: S-0591
owner: alex
created: 2026-11-07T21:48:21Z
updated: 2026-11-07T21:48:21Z
transitions: []
stream: S-0591
tags: [issues]
touches: [design/issues/I-0303-testgatewaylitellm-fails-its-setup-check-the-template-s-settings-run-flai-guard.md, design/issues/summary.md]
after: [T-2592]
---
# T-2593 Close I-0303 with the reason naming the guardSettings fix and its test

## Work

Run `flai issue close I-0303 --reason "<reason>"` in the story's worktree, the reason naming what fixed it: `guardSettings` counts and rewrites only the hooks' invocations of `flai guard`, not the words of the S-0518 message, and the new test named in T-2592 checks it on the template's settings under `-short`. Check that `design/issues/summary.md` no longer lists I-0303.

Waits for T-2592, whose fix and test name the reason.

## Done when

- I-0303's status is closed, with that reason.
- `flai check --strict` is clean on the issue and the summary.

## Notes
