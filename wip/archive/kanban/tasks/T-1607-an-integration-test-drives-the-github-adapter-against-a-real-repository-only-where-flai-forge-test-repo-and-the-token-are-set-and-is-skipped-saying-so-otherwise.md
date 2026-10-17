---
id: T-1607
type: task
nature: feature
title: An integration test drives the GitHub adapter against a real repository only where FLAI_FORGE_TEST_REPO and the token are set, and is skipped saying so otherwise
status: done
parent: S-0392
owner: alex
created: 2026-10-16T04:33:48Z
updated: 2026-10-17T07:12:49Z
transitions:
  - to: ready
    at: 2026-10-17T07:03:23Z
    by: agent-S-0392
  - to: in-progress
    at: 2026-10-17T07:03:23Z
    by: agent-S-0392
  - to: done
    at: 2026-10-17T07:12:49Z
    by: agent-S-0392
stream: S-0392
tags: [cli]
touches: [flai/tests/integration/forge_github_test.go]
after: [T-1603]
usage:
  source: log
  seconds: 566
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 5529
      cache_read: 681159
      cache_write: 23240
      cost: 0.3796
---
# T-1607 An integration test drives the GitHub adapter against a real repository only where FLAI_FORGE_TEST_REPO and the token are set, and is skipped saying so otherwise

## Work

- Add `flai/tests/integration/forge_github_test.go`. It skips, with `t.Skip` naming the variable that is missing, unless both `FLAI_FORGE_TEST_REPO` (`owner/name`) and the token variable are set, as `mcp_http_test.go` skips without its URL.
- Where both are set, it does the following on that repository, cleaning up after itself:
  - pushes a scratch branch
  - opens a pull request, finds it by its head branch, and updates its body
  - reads its status
  - closes it
- It never merges into the test repository's default branch.
- Waits for T-1603: it drives the adapter.

## Done when

Without the variables the test is skipped with a message naming them, and `scripts/integration.sh` passes. With them it opens, finds, updates, reads, and closes a pull request.

## Notes
