---
id: D-0533
title: "The gateway tests' guardSettings counts and rewrites only the template hooks' runs of flai guard, after $( or a semicolon and a space, never the hooks' message 'flai guard did not run', so a hook change keeps that shape"
kind: decision
sources: [S-0591, I-0303]
topics: [testing, flai, go, code, harness]
paths: [flai/internal/serve/gateway_test.go, template/root/.claude/settings.json]
date: 2026-11-09
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0533 The gateway tests' guardSettings counts and rewrites only the template hooks' runs of flai guard, after $( or a semicolon and a space, never the hooks' message 'flai guard did not run', so a hook change keeps that shape

I-0303: `guardSettings` in `flai/internal/serve/gateway_test.go` counted every `flai guard` substring in `template/root/.claude/settings.json` and wanted 4. S-0518 added the message `flai guard did not run` to the two refusing hooks, which made 6, so `TestGatewayLiteLLM` failed its setup check.

S-0591 fixed the test, not the settings, because the message is the guard hook's own error text. Its Decisions (`wip/archive/agents/S-0591.md`) say:

```text
A run of `flai guard` is matched as a command, after `$(` or after a semicolon and a space; the message is inside single quotes after `printf '`, so it never matches. `TestGuardSettingsRunsTheBuiltFlai` checks four rewritten runs, two untouched messages, and valid JSON, and is not named `TestGateway*`, so the `go-test` tier runs it under `-short`.
```

So a story that changes the template's guard hooks keeps each run of `flai guard` after `$(` or `;`, and any text that names flai guard inside quotes. When it adds or removes a run, it changes the count `TestGuardSettingsRunsTheBuiltFlai` expects. That test runs under `flai test`; the paid `TestGateway*` checks do not (D-0451).

Sources: S-0591, I-0303.
