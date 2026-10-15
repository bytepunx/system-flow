---
id: T-1543
type: task
nature: feature
title: An edit model holds changed tokens apart from the saved ones, renders them as flai does for the preview, and writes them back into the DTCG token file, leaving the rest as it was
status: done
parent: S-0383
owner: alex
created: 2026-10-13T17:47:27Z
updated: 2026-10-15T02:18:13Z
transitions:
  - to: ready
    at: 2026-10-15T01:58:59Z
    by: agent-S-0383
  - to: in-progress
    at: 2026-10-15T01:58:59Z
    by: agent-S-0383
  - to: done
    at: 2026-10-15T02:18:13Z
    by: agent-S-0383
stream: S-0383
tags: [dashboard]
touches: [flaiover/src/lib/design/tokens.ts, flaiover/src/lib/design/tokens.test.ts]
usage:
  source: log
  seconds: 1154
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 69
      output: 47870
      cache_read: 4360146
      cache_write: 170773
      cost: 2.7706
---
# T-1543 An edit model holds changed tokens apart from the saved ones, renders them as flai does for the preview, and writes them back into the DTCG token file, leaving the rest as it was

## Work

`flaiover/src/lib/design/tokens.ts` is the edit model of the Design page:

- It holds the saved tokens per theme, as `design.get` answers them (`DesignTheme[]`), and the edits keyed by theme and token path.
- It renders an edited token's `value`, `css`, and `variables` as flai's `flai/internal/designsystem/css.go` does, for the kinds the editor changes: colour, dimension, font family, and a typography step's font size and line height. An alias of an edited token follows it.
- With the raw token file, a light edit of a token the dark context does not override changes the dark theme too, as the resolver does.
- It checks a value as flai's `design.tokens` does, before saving: a colour in `srgb` with its hex or in `oklch`, a dimension in `px` or `rem`, and a non-empty font family.
- It says which tokens changed, discards back to the saved ones, and composes the token file's new raw text. Each edited `$value` is replaced where it stands, and a dark edit with no override is added to the dark context. Every other byte is kept, including key order, formatting, and unknown keys.

It waits for nothing. It was retargeted from the front matter at the start (TH-0454): the tokens are the DTCG file of ADR-0146.

## Done when

- `tokens.test.ts` passes. It covers each kind of value, the alias, the dark inheritance, a bad value, discard, and a round trip with no edits that returns the file byte for byte.

## Notes

- Layer 1, beside T-1544.
