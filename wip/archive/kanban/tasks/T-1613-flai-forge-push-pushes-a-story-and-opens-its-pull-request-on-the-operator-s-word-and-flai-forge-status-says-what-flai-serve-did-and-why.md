---
id: T-1613
type: task
nature: feature
title: flai forge push pushes a story and opens its pull request on the operator's word, and flai forge status says what flai serve did and why
status: done
parent: S-0393
owner: alex
created: 2026-10-16T04:35:42Z
updated: 2026-10-17T08:33:38Z
transitions:
  - to: ready
    at: 2026-10-17T08:16:10Z
    by: agent-S-0393
  - to: in-progress
    at: 2026-10-17T08:16:14Z
    by: agent-S-0393
  - to: done
    at: 2026-10-17T08:33:38Z
    by: agent-S-0393
stream: S-0393
tags: [cli]
touches: [flai/cmd/forge.go, flai/cmd/forge_test.go]
after: [T-1612]
usage:
  source: log
  seconds: 1044
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 79
      output: 33296
      cache_read: 4286776
      cache_write: 126299
      cost: 2.2441
---
# T-1613 flai forge push pushes a story and opens its pull request on the operator's word, and flai forge status says what flai serve did and why

## Work

- Add `flai forge push S-nnnn` to `flai/cmd/forge.go`. It calls T-1612's function for one story, with the same checks, and prints the pull request's URL.
- Add `flai forge status S-nnnn`. It prints:
  - the pull request's number and URL
  - whether the remote branch matches the local head
  - when nothing was pushed, why: the mode, the host action, the token's variable, or the last failure
- `--json` gives the same fields for both commands.
- Test both against the fake forge.
- Waits for T-1612: both commands call its code.

## Done when

Both commands work against the fake forge in every case T-1612 covers, and their tests pass.

## Notes
