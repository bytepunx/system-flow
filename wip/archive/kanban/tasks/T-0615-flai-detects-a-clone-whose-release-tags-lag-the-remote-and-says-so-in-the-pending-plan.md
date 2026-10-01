---
id: T-0615
type: task
nature: remediation
title: flai detects a clone whose release tags lag the remote and says so in the pending plan
status: done
parent: S-0174
owner: arobson
created: 2026-10-01T07:51:03Z
updated: 2026-10-01T07:59:57Z
transitions:
  - to: ready
    at: 2026-10-01T07:51:10Z
    by: agent-S-0174
  - to: in-progress
    at: 2026-10-01T07:51:10Z
    by: agent-S-0174
  - to: done
    at: 2026-10-01T07:59:57Z
    by: agent-S-0174
stream: S-0174
tags: []
touches: [flai/internal/release, flai/internal/preview, flai/cmd/release.go]
usage:
  source: log
  seconds: 527
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 73
      output: 24370
      cache_read: 6448965
      cache_write: 71034
      cost: 2.3458
---
# T-0615 flai detects a clone whose release tags lag the remote and says so in the pending plan

## Work

- In `flai/internal/release`, ask the remote the branch tracks (else `origin`) for its release tags with one `git ls-remote --tags --refs`, compare the highest `<name>/vX.Y.Z` per code component with the highest local one, and say which components lag, or that the remote could not be asked, with the fix (`git fetch --tags <remote>`). Bounded by a timeout; the answer is kept a short while per root so the board asking again does not ask the network each time.
- `preview.Publish` (`publish.preview`, `flai release --pending --dry-run`) carries the check: with a component behind it offers no plan and says what fixes it; with the remote unreachable it keeps the plan and warns.

## Done when

- Behaviour tests cover a clone whose tags lag the remote, one in step with it, and an unreachable remote; `make test` and lint pass.

## Notes
