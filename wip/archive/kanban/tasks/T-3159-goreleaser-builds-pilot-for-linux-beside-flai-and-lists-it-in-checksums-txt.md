---
id: T-3159
type: task
nature: feature
title: GoReleaser builds pilot for Linux beside flai and lists it in checksums.txt
status: done
parent: S-0734
owner: alex
created: 2026-11-12T15:57:26Z
updated: 2026-11-13T08:02:37Z
transitions:
  - to: ready
    at: 2026-11-13T07:56:23Z
    by: agent-S-0734
  - to: in-progress
    at: 2026-11-13T07:56:23Z
    by: agent-S-0734
  - to: done
    at: 2026-11-13T08:02:37Z
    by: agent-S-0734
stream: S-0734
tags: [cli]
touches: [flai/.goreleaser.yaml]
usage:
  source: log
  seconds: 374
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 16
      output: 5899
      cache_read: 1365585
      cache_write: 37637
      cost: 0.6572
---
# T-3159 GoReleaser builds pilot for Linux beside flai and lists it in checksums.txt

## Work

Add a build of `./pilot` to `flai/.goreleaser.yaml` for `linux/amd64` and `linux/arm64`, static (`CGO_ENABLED=0`), stamped as flai is. Archive it, or ship it as a binary, as S-0733's ADR names, so it lands in `checksums.txt`, which the release key signs. Keep the version template's `trimprefix .Tag "flai/v"`. Check the result with `goreleaser release --snapshot --clean --skip=publish` in the layer image, where GoReleaser is, and list what it wrote.

It waits for no task: GoReleaser's file is its own.

## Done when

- A snapshot build writes `pilot` for both Linux architectures, and `checksums.txt` names them.
- `flai`'s own archives are unchanged.

## Notes
