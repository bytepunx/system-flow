---
id: T-1612
type: task
nature: feature
title: flai serve pushes each story in review under pull-request acceptance and opens or updates its pull request, raising a failure on one thread on the story
status: done
parent: S-0393
owner: alex
created: 2026-10-16T04:35:39Z
updated: 2026-10-17T08:16:03Z
transitions:
  - to: ready
    at: 2026-10-17T07:58:15Z
    by: agent-S-0393
  - to: in-progress
    at: 2026-10-17T07:58:17Z
    by: agent-S-0393
  - to: done
    at: 2026-10-17T08:16:03Z
    by: agent-S-0393
stream: S-0393
tags: [cli]
touches: [flai/internal/serve/forge.go, flai/internal/serve/forge_test.go, flai/internal/serve/serve.go, flai/internal/forge/forge.go, flai/cmd/forge.go, flai/cmd/serve.go]
after: [T-1609, T-1610]
usage:
  source: log
  seconds: 1066
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 110
      output: 46388
      cache_read: 5972335
      cache_write: 175959
      cost: 3.1265
---
# T-1612 flai serve pushes each story in review under pull-request acceptance and opens or updates its pull request, raising a failure on one thread on the story

## Work

- In `flai/internal/serve/forge.go`, run on flai serve's minute loop, wired in `serve.go`, for each served project with `acceptance.mode: pull-request` and the `forge` host action on. For each story in review:
  - When `story/S-nnnn` on `origin` is missing or differs from the local head, push it with the host's git. Use `--force-with-lease=<branch>:<last pushed head>` when the remote head is not an ancestor, and never force any other branch.
  - `Find` its pull request. `Open` one with T-1609's title and body when there is none, else `Update` the body when it changed.
  - Keep the last pushed head per story under `.flai-cache`, so that the lease names it.
- On any failure, open one thread on the story with `threads.New`, as `upgrade.go` does. The thread names the reason and what clears it. Resolve the thread once a later look succeeds.
- Expose the function, so that `flai forge push` calls the same code.
- Test it with a scratch remote and T-1603's fake forge: first push, update, rebase with lease, a lease refused, a failure thread opened and then resolved, the host action off, and the mode `local`.
- Waits for T-1609 (the body) and T-1610 (the host action).

## Done when

A story moved to review in a scratch project gets its branch pushed and one pull request, a second move to review updates the same one, a failure opens and later resolves one thread, and the tests pass.

## Notes
