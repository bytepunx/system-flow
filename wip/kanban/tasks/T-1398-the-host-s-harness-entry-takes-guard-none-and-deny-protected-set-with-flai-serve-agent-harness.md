---
id: T-1398
type: task
nature: feature
title: The host's harness entry takes guard none and deny_protected, set with flai serve agent harness
status: backlog
parent: S-0355
owner: alex
created: 2026-10-08T08:49:02Z
updated: 2026-10-08T08:54:19Z
transitions: []
stream: S-0355
tags: [cli]
touches: [flai/internal/config/config.go, flai/internal/config/config_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go, docs/operators/settings.md, docs/users/flai-reference.md]
---
# T-1398 The host's harness entry takes guard none and deny_protected, set with flai serve agent harness

## Work

- `config.HarnessHost` gains `Guard string` (`json:"guard,omitempty"`, only `none` accepted) and `DenyProtected bool` (`json:"deny_protected,omitempty"`).
- `flai serve agent harness <name> --guard none|""` and `--deny-protected=true|false`, beside `--program` and the arguments; the help says what each gives up.
- Rows in `docs/operators/settings.md` for both keys and both flags.

First layer.

## Done when

- `config_test.go` covers a round trip and a refused value; `serve_actions_test.go` sets and clears both.
- `flai test flai/internal/config/ flai/cmd/ docs/operators/settings.md` passes.

## Notes

Layer 1 of S-0355.
