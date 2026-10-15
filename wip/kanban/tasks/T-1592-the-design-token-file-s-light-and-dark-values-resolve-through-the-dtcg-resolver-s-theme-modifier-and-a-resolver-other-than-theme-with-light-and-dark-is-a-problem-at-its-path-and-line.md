---
id: T-1592
type: task
nature: feature
title: The design token file's light and dark values resolve through the DTCG Resolver's theme modifier, and a resolver other than theme with light and dark is a problem at its path and line
status: done
parent: S-0381
owner: alex
created: 2026-10-14T23:53:09Z
updated: 2026-10-15T00:17:48Z
transitions:
  - to: ready
    at: 2026-10-15T00:09:52Z
    by: agent-S-0381
  - to: in-progress
    at: 2026-10-15T00:09:52Z
    by: agent-S-0381
  - to: done
    at: 2026-10-15T00:17:47Z
    by: agent-S-0381
stream: S-0381
tags: [cli]
touches: [flai/internal/designsystem/resolver.go, flai/internal/designsystem/resolver_test.go, flai/internal/designsystem/tokens.go]
after: [T-1532]
usage:
  source: log
  seconds: 475
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 28
      output: 123
      cache_read: 992480
      cache_write: 104808
      cost: 0.5067
---
# T-1592 The design token file's light and dark values resolve through the DTCG Resolver's theme modifier, and a resolver other than theme with light and dark is a problem at its path and line

## Work

In `flai/internal/designsystem/resolver.go`, implement the subset of the DTCG 2025.10 Resolver module ADR-0146 names: one modifier, `theme`, with the contexts `light` and `dark`, and nothing more. Settle against the spec whether the dark values sit in the token file the manifest names or in a set and a resolver document beside it, and record the layout in the package comment for T-1537 to describe; the starter T-1535 writes follows it. `Resolve` answers the tokens of each theme, a dark value falling back to the light one where the token has one value for both. `Load` reads the token file at the path it is given, runs T-1532's `Validate` and the resolver's own checks, and answers the themes and every problem, each with its token path and line: a modifier other than `theme`, a context other than `light` and `dark`, a context that names a token the file lacks, and a reference that does not resolve. Table tests cover each problem, a file with one value for both themes, and one with dark values. It waits for T-1532, whose token types it resolves.

## Done when

- `go test ./internal/designsystem/` passes with a test per resolver problem and both valid layouts.
- The package comment says where the dark values sit.

## Notes

- Layer 2, beside the colour task, with which it shares no file. Written by planner-S-0381.
