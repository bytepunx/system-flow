---
id: T-2231
type: task
nature: remediation
title: I-0217 is closed with flai issue close saying what fixed it, and the story's criteria are ticked
status: backlog
parent: S-0519
owner: alex
created: 2026-10-28T13:45:13Z
updated: 2026-10-28T13:45:41Z
transitions: []
stream: S-0519
tags: [cli]
touches: [design/issues/I-0217-permission-prompt-over-the-http-mcp-server-outlives-claude-code-s-tool-call-timeout-so-an-agent-in-a-container-sees-its-protected-write-time-out-while-the-thread-to-the-operator-stays-open.md, design/issues/summary.md]
after: [T-2229]
---
# T-2231 I-0217 is closed with flai issue close saying what fixed it, and the story's criteria are ticked

## Work

From the story's worktree, since issue commands write under the checkout they run in, run `flai issue close I-0217 --reason "..."` naming the ADR and the bound over HTTP that the code task built, once its test passes. Then tick criteria 1 and 2 with `flai criteria tick S-0519 <n>`. It waits for the code task: the issue is closed by the fix, not before it.

## Done when

I-0217 is `closed` with a reason naming the fix, `design/issues/summary.md` no longer lists it as open, both criteria of S-0519 are ticked, and `flai check --strict` passes.

## Notes

The measurement task may have filed an issue for another long call the timeout cuts short; that issue stays open and is not this story's to close.
