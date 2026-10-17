---
id: T-1606
type: task
nature: feature
title: flai forge check says whether the token the host names can read and push to the repository behind origin
status: done
parent: S-0392
owner: alex
created: 2026-10-16T04:33:44Z
updated: 2026-10-17T07:12:46Z
transitions:
  - to: ready
    at: 2026-10-17T07:03:22Z
    by: agent-S-0392
  - to: in-progress
    at: 2026-10-17T07:03:22Z
    by: agent-S-0392
  - to: done
    at: 2026-10-17T07:12:46Z
    by: agent-S-0392
stream: S-0392
tags: [cli]
touches: [flai/cmd/forge.go, flai/cmd/forge_test.go, flai/cmd/root.go]
after: [T-1602, T-1603, T-1604]
usage:
  source: log
  seconds: 564
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 33
      output: 14653
      cache_read: 1805176
      cache_write: 61589
      cost: 1.006
---
# T-1606 flai forge check says whether the token the host names can read and push to the repository behind origin

## Work

- Add `flai forge check` in `flai/cmd/forge.go`, and register the `forge` group in `flai/cmd/root.go`. It reads the project's `acceptance.forge`, the `origin` remote's URL, and the host's forge entry. It resolves the token, calls `Access`, and prints:
  - the forge
  - the repository
  - the variable the token is read from
  - read: yes or no
  - push: yes or no
- `--json` gives the same fields. The command exits non-zero, with the reason, when the variable is unset, the remote is not on the forge, or access fails.
- It never prints the token.
- Test it against the fake server from T-1603, with the variable set and unset.
- Waits for T-1602 (the manifest keys), T-1603 (the contract), and T-1604 (the host's variable): it uses all three.

## Done when

`flai forge check` answers for a project with and without the token set, exits non-zero with the reason when it cannot reach the forge, and its tests pass.

## Notes
