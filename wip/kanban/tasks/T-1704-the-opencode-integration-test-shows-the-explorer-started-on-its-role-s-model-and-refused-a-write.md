---
id: T-1704
type: task
nature: feature
title: The OpenCode integration test shows the explorer started on its role's model and refused a write
status: backlog
parent: S-0409
owner: alex
created: 2026-10-16T05:17:41Z
updated: 2026-10-16T05:17:41Z
transitions: []
stream: S-0409
tags: [cli]
touches: [flai/tests/integration/opencode_test.go]
after: [T-1702, T-1703]
---
# T-1704 The OpenCode integration test shows the explorer started on its role's model and refused a write

## Work

- Extend `flai/tests/integration/opencode_test.go`: a story with `roles.explore` naming a model hands a search to the explorer; assert from the JSON output that the explorer ran on that model and that its file write was refused. Where `Guard` is false, assert the start is refused naming `agent.harnesses.opencode.guard`.
- Keep the skip on a missing or unknown `opencode`.
- Waits for T-1702 (the agents in the configuration) and T-1703 (the read-only mark).

## Done when

The test passes on the known `opencode` version and is skipped, saying so, without it.

## Notes
