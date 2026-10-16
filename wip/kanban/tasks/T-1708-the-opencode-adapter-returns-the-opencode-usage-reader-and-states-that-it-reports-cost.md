---
id: T-1708
type: task
nature: feature
title: The opencode adapter returns the OpenCode usage reader and states that it reports cost
status: backlog
parent: S-0410
owner: alex
created: 2026-10-16T05:18:51Z
updated: 2026-10-16T05:23:22Z
transitions: []
stream: S-0410
tags: [cli]
touches: [flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/harness_test.go]
after: [T-1706, T-1707]
---
# T-1708 The opencode adapter returns the OpenCode usage reader and states that it reports cost

## Work

- In `opencode.go`, return `usage.OpenCode()` from `Reader` in place of the unread reader, set the cost capability, and pass the provider entry T-1707 writes into the start.
- Update `TestEachAdapterStatesItsCapabilities` and the adapter's start tests with and without a provider.
- Waits for T-1706 (the reader) and T-1707 (the provider entry, in the configuration this start writes).

## Done when

The adapter gives the reader, states the capability, starts through a provider, and the harness tests pass.

## Notes
