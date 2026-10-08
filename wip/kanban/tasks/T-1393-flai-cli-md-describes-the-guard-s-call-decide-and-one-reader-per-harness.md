---
id: T-1393
type: task
nature: improvement
title: flai-cli.md describes the guard's Call, Decide, and one reader per harness
status: backlog
parent: S-0353
owner: alex
created: 2026-10-08T08:47:06Z
updated: 2026-10-08T08:47:06Z
transitions: []
stream: S-0353
tags: [cli]
touches: [design/system/flai-cli.md]
after: [T-1392]
---
# T-1393 flai-cli.md describes the guard's Call, Decide, and one reader per harness

## Work

- `design/system/flai-cli.md`, where `flai guard` and the internal structure are described: the neutral `Call`, the three verdicts, `Decide`, the reader per harness with Claude Code's as the one built, `--harness`, and that the rules and their texts did not change. Cite ADR-0130.

It waits for T-1392, so that it describes the guard as built.

## Done when

- The section matches the code.
- `flai test design/system/flai-cli.md` passes.

## Notes

Layer 3 of S-0353.
