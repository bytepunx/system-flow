---
id: T-1031
type: task
nature: feature
title: The host API reports and checks the shared paths, and settings.shared adds or removes them through flai shared
status: backlog
parent: S-0295
owner: alex
created: 2026-10-06T12:16:17Z
updated: 2026-10-06T12:16:36Z
transitions: []
stream: S-0295
tags: [flai]
touches: [flai/internal/hostapi/settings.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/contract_test.go, flai/cmd/serve_actions.go, flai/cmd/serve_actions_test.go]
after: [T-1028]
---
# T-1031 The host API reports and checks the shared paths, and settings.shared adds or removes them through flai shared

## Work

Expose the shared paths over the host API, as the other project settings are. The operator asked on TH-0180 for checking as well as managing.

- **Read.** `settings.get` reports the project's patterns, through `hostSettings` in `flai/cmd/serve_actions.go`, so the dashboard can show them.
- **Check.** A read method (proposed `settings.shared_check`) takes paths or a story ID and answers what `flai shared check --json` answers: for each entry, whether it is shared and which pattern matched. It changes nothing, so it is gated like `settings.get`.
- **Write.** A new `settings.shared` method in `flai/internal/hostapi/settings.go` adds or removes one pattern, gated by the `settings` host action like `settings.default_agent`. It delegates to the `flai shared add|remove` command T-1028 adds (ADR-0016) and answers with what changed, or with flai's refusal of an invalid or duplicate pattern as a channel error.

Add both methods to the contract test, if `contract_test.go` lists the methods.

This task waits for T-1028, whose command it runs.

## Done when

- `writes_test.go`, or the settings tests beside it, covers add, remove, a refused pattern, and the refusal with the settings action off.
- A test covers the check method on a shared path, a path that is not shared, and a story ID.
- `serve_actions_test.go` shows `settings.get` carrying the list.
- `scripts/flai-test.sh` passes.

## Notes
