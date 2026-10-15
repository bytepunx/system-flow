---
id: T-1534
type: task
nature: feature
title: flai check reports a design system token that does not parse as design.tokens with its key and line
status: done
parent: S-0381
owner: alex
created: 2026-10-13T17:46:35Z
updated: 2026-10-15T00:29:28Z
transitions:
  - to: ready
    at: 2026-10-15T00:22:33Z
    by: agent-S-0381
  - to: in-progress
    at: 2026-10-15T00:22:33Z
    by: agent-S-0381
  - to: done
    at: 2026-10-15T00:29:28Z
    by: agent-S-0381
stream: S-0381
tags: [cli]
touches: [flai/internal/check/design.go, flai/internal/check/design_test.go, flai/internal/check/check.go]
after: [T-1590, T-1592]
usage:
  source: log
  seconds: 415
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 34
      output: 17599
      cache_read: 1653153
      cache_write: 79921
      cost: 1.1217
---
# T-1534 flai check reports a design system token that does not parse as design.tokens with its key and line

## Work

Add the rule `design.tokens` in `flai/internal/check/design.go`, called from `Run` in `check.go`, as `experiments` and `analysis` are: it calls `designsystem.Load` on the path the manifest's `design.tokens` names (T-1590) and reports one error per problem, on the token file's path and the token's line, naming the token's path. A project with no token file has no finding. Contrast is never a finding (ADR-0148). Test a valid file, one with a token problem, one with a resolver problem, and none. It waits for T-1590, whose key names the path, and T-1592, whose `Load` it reports.

## Done when

- `go test ./internal/check/` passes with the four cases.
- `flai check --strict` on this repository is unchanged.

## Notes

- Layer 3, beside T-1533. Rewritten by planner-S-0381 from the front matter plan, as TH-0453 approved.
