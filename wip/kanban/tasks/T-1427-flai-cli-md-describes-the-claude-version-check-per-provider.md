---
id: T-1427
type: task
nature: improvement
title: flai-cli.md describes the claude version check per provider
status: backlog
parent: S-0360
owner: alex
created: 2026-10-08T08:59:51Z
updated: 2026-10-08T08:59:51Z
transitions: []
stream: S-0360
tags: [cli]
touches: [design/system/flai-cli.md]
after: [T-1426]
---
# T-1427 flai-cli.md describes the claude version check per provider

## Work

- `design/system/flai-cli.md`, where the check of each new `claude` version is described: it runs through the project's provider, the record per version and provider, and the skip on an unset key.

It waits for T-1426, so that it describes the check as built.

## Done when

- `flai test design/system/flai-cli.md` passes.

## Notes

Layer 2 of S-0360.
