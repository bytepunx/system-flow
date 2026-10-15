---
id: T-1532
type: task
nature: feature
title: flai/internal/designsystem reads the DTCG token file and the guidance's headings, and reports each token it does not read at its path and line
status: done
parent: S-0381
owner: alex
created: 2026-10-13T17:46:26Z
updated: 2026-10-15T00:09:25Z
transitions:
  - to: ready
    at: 2026-10-14T23:57:50Z
    by: agent-S-0381
  - to: in-progress
    at: 2026-10-14T23:57:51Z
    by: agent-S-0381
  - to: done
    at: 2026-10-15T00:09:25Z
    by: agent-S-0381
stream: S-0381
tags: [cli]
touches: [flai/internal/designsystem/tokens.go, flai/internal/designsystem/tokens_test.go, flai/internal/designsystem/testdata/full.tokens.json]
usage:
  source: log
  seconds: 694
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 42
      output: 21723
      cache_read: 2040491
      cache_write: 98647
      cost: 1.3845
    - model: claude-sonnet-5-5
      input: 44
      output: 12003
      cache_read: 1390155
      cache_write: 103370
      cost: 0.5176
---
# T-1532 flai/internal/designsystem reads the DTCG token file and the guidance's headings, and reports each token it does not read at its path and line

## Work

Add the package `flai/internal/designsystem`. `Read` parses the DTCG 2025.10 token file at the path it is given (ADR-0146), keeping each token's line, since `encoding/json` gives none: decode with offsets, or walk the tokens and map each offset to its line. It reads the types `color`, `fontFamily`, `dimension`, `shadow`, and `typography`, with a group's `$type` inherited by its tokens, in the groups `color`, `font`, `text`, `spacing`, `radius`, and `shadow`, as `design/system/design-system-tool.md` § The token schema lays out. It also lists the headings of the guidance, `design/system/design-system.md`. `Validate` returns one problem per token it does not read (ADR-0148), each with its token path and line: a file that does not parse, a type flai does not read, a colour space other than `srgb` with its six-digit `hex` or `oklch`, alpha outside shadows and extra tokens, a length in a unit other than `rem` or `px`, and a missing colour role, font, type step, spacing, radius, or shadow. Any colour name outside the roles is an extra token, read and validated. Table tests cover each problem and a valid file. It waits for nothing.

## Done when

- `go test ./internal/designsystem/` passes with a test per problem, each asserting its path and line, and a valid file.

## Notes

- Layer 1, beside T-1590 and T-1591. The path is the caller's: T-1590's manifest key answers it. Rewritten by planner-S-0381 from the front matter plan, as TH-0453 approved.
