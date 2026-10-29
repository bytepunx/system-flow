---
id: T-2329
type: task
nature: remediation
title: Close I-0227 with flai issue close, naming the timeout finding, the tier's timeout, and the faster prime tests
status: backlog
parent: S-0540
owner: alex
created: 2026-10-29T07:37:46Z
updated: 2026-10-29T07:37:46Z
transitions: []
stream: S-0540
tags: [flai, issues]
touches: [design/issues/I-0227-flai-test-s-go-test-tier-runs-the-mcpserver-package-with-race-into-go-test-s-10-minute-default-timeout-on-a-loaded-host-and-names-only-the-package-so-a-story-that-changes-mcpserver-cannot-pass-its-tests.md, design/issues/summary.md]
after: [T-2326, T-2327, T-2328]
---
# T-2329 Close I-0227 with flai issue close, naming the timeout finding, the tier's timeout, and the faster prime tests

## Work

Criterion 2. In the story's worktree, run `flai issue close I-0227 --reason "<reason>"`. The reason names the three changes: the `go-test-json` format names the tests running at a timeout (T-2326), the go-test tier runs with `-timeout 20m` (T-2327), and the two prime-on-this-repository tests take less time under `-race` (T-2328), with the before and after figures T-2328 recorded. flai writes the issue file and `design/issues/summary.md`.

Before you close it, run `flai test flai/internal/mcpserver` once in the worktree and record its time and result in the narrative. That shows the cause no longer occurs on this host, for criterion 1.

It waits for T-2326, T-2327, and T-2328: the reason names what each changed, and the closing run needs all three.

## Done when

- I-0227's status is closed, its reason names the three changes, and `design/issues/summary.md` no longer lists it.
- `flai test flai/internal/mcpserver` passed in the worktree, and its time is in the narrative.
- `flai check --strict` passes.

## Notes
