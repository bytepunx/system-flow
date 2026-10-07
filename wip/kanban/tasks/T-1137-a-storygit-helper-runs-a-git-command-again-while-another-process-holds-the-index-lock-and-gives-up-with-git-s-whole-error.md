---
id: T-1137
type: task
nature: remediation
title: A storygit helper runs a git command again while another process holds the index lock, and gives up with git's whole error
status: in-progress
parent: S-0307
owner: alex
created: 2026-10-07T01:12:33Z
updated: 2026-10-07T01:42:24Z
transitions:
  - to: ready
    at: 2026-10-07T01:42:24Z
    by: agent-S-0307
  - to: in-progress
    at: 2026-10-07T01:42:24Z
    by: agent-S-0307
stream: S-0307
tags: [flai]
touches: [flai/internal/storygit/indexlock.go, flai/internal/storygit/indexlock_test.go]
---
# T-1137 A storygit helper runs a git command again while another process holds the index lock, and gives up with git's whole error

## Work

Add `flai/internal/storygit/indexlock.go`: a function that runs a git command through an `execx.Runner` and, when git's output says another process holds the index lock (`Unable to create '…/index.lock': File exists`), waits and runs it again. Wait a short time that grows, for a total of about ten seconds, since the other writer (flai serve's replans, a planner's or a story's autocommit, another acceptance) holds the lock for well under a second. Any other failure returns at once. When the lock is still held after the last try, return git's whole output with a line saying the index lock was held throughout. Never delete the lock file: it may belong to a live process.

Take the wait as a parameter, or a package variable a test can shorten, so tests do not sleep for seconds.

The file is new and kept apart from `flai/internal/storygit/commit.go`, which S-0275 plans, so the two stories do not change one file. It waits for no task: nothing else in S-0307 is needed to write it.

## Done when

- `indexlock_test.go` covers, with a stand-in runner: a command that succeeds at once runs once; a command that fails with the index-lock message and then succeeds is run again and succeeds; one that keeps failing with it gives up after the last try with git's output and the held-lock line; any other failure returns at once without a retry.
- A test against a real repository holds `.git/index.lock`, removes it from a goroutine after a short wait, and the commit through the helper succeeds; it skips when git is not installed.
- `flai test flai/internal/storygit` passes.

## Notes
