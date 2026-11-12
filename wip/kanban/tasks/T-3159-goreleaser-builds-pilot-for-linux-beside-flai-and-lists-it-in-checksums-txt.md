---
id: T-3159
type: task
nature: feature
title: GoReleaser builds pilot for Linux beside flai and lists it in checksums.txt
status: backlog
parent: S-0734
owner: alex
created: 2026-11-12T15:57:26Z
updated: 2026-11-12T15:57:26Z
transitions: []
stream: S-0734
tags: [cli]
touches: [flai/.goreleaser.yaml]
---
# T-3159 GoReleaser builds pilot for Linux beside flai and lists it in checksums.txt

## Work

Add a build of `./pilot` to `flai/.goreleaser.yaml` for `linux/amd64` and `linux/arm64`, static (`CGO_ENABLED=0`), stamped as flai is. Archive it, or ship it as a binary, as S-0733's ADR names, so it lands in `checksums.txt`, which the release key signs. Keep the version template's `trimprefix .Tag "flai/v"`. Check the result with `goreleaser release --snapshot --clean --skip=publish` in the layer image, where GoReleaser is, and list what it wrote.

It waits for no task: GoReleaser's file is its own.

## Done when

- A snapshot build writes `pilot` for both Linux architectures, and `checksums.txt` names them.
- `flai`'s own archives are unchanged.

## Notes
