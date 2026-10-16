---
id: T-1613
type: task
nature: feature
title: flai forge push pushes a story and opens its pull request on the operator's word, and flai forge status says what flai serve did and why
status: backlog
parent: S-0393
owner: alex
created: 2026-10-16T04:35:42Z
updated: 2026-10-16T04:35:42Z
transitions: []
stream: S-0393
tags: [cli]
touches: [flai/cmd/forge.go, flai/cmd/forge_test.go]
after: [T-1612]
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
