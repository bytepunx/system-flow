---
id: T-2214
type: task
nature: remediation
title: "harness.None is a reserved name: Elsewhere reports an agent that runs elsewhere, and harness.For refuses it saying flai serve starts nothing for it"
status: backlog
parent: S-0516
owner: alex
created: 2026-10-28T13:13:27Z
updated: 2026-10-28T13:13:27Z
transitions: []
stream: S-0516
tags: [cli, serve]
touches: [flai/internal/harness/harness.go, flai/internal/harness/harness_test.go]
after: [T-2213]
---
# T-2214 harness.None is a reserved name: Elsewhere reports an agent that runs elsewhere, and harness.For refuses it saying flai serve starts nothing for it

## Work

In `flai/internal/harness/harness.go`, beside `Command`, add the constant `None = "none"` and a helper, `Elsewhere(a *manifest.Agent) bool`, true when the agent names `None`, that every start path in `flai serve` checks before it reaches `For`. Make `For` refuse `None` with one error, worded for the operator: the story's agent runs elsewhere (`harness: none`), so `flai serve` starts nothing for it; start it yourself with `flai story start`, or name a harness it can start. Keep `None` out of the `Adapters` map and out of the list the unknown-harness error prints (`it can start claude-code, command`), since nothing can start it. Leave the manifest's shape check alone: `none` matches `harnessPattern` already. Waits for T-2213, which fixes the name and the words the ADR records.

## Done when

- `harness.None` and `harness.Elsewhere` exist, `Elsewhere` is false for a nil agent, an empty harness, and every other name
- `harness.For` on an agent with harness `none` returns an error naming `flai story start` and no adapter, whether or not a command is set on the host, and a test in `harness_test.go` checks the error's text beside the existing "names no harness" and "cannot start harness" cases
- `TestEachAdapterStatesItsCapabilities` still passes: `None` is not an adapter

## Notes

Layer 2, with T-2213 before it. S-0502, in progress, touches `harness.go` and `harness_test.go` in T-2145: the story is held in ready until S-0502 moves to review, and the first sync after it lands brings its change.
