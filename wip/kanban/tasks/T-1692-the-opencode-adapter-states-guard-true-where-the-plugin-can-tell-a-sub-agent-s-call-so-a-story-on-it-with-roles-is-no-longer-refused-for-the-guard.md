---
id: T-1692
type: task
nature: feature
title: The opencode adapter states Guard true where the plugin can tell a sub-agent's call, so a story on it with roles is no longer refused for the guard
status: backlog
parent: S-0407
owner: alex
created: 2026-10-16T05:15:38Z
updated: 2026-10-16T05:27:12Z
transitions: []
stream: S-0407
tags: [cli]
touches: [flai/internal/harness/opencode.go, flai/internal/harness/opencode_test.go, flai/internal/harness/harness_test.go, flai/internal/serve/agents_test.go]
after: [T-1690]
---
# T-1692 The opencode adapter states Guard true where the plugin can tell a sub-agent's call, so a story on it with roles is no longer refused for the guard

## Work

- Where S-0406 found a sub-agent sign, set `Guard: true` in the adapter's `Capabilities` in `flai/internal/harness/opencode.go`, and update `TestEachAdapterStatesItsCapabilities` and the adapter's tests.
- Where it found none, leave `Guard` false and say so in a comment naming the finding. Roles stay refused until role models arrive in the roles story either way.
- In `flai/internal/serve/agents_test.go`, test `startAllowed` with the opencode adapter's capabilities and a story whose agent has roles. Where `Guard` is true, the guard's refusal is gone. Where it is false, the refusal stands and names `--guard none`. This shows criterion 3 in a test (added on TH-0470, as the orchestrator approved).
- Waits for T-1690: the capability is true only once the reader judges sub-agents.

## Done when

The capability matches the finding, and the harness and serve tests pass.

## Notes
