---
id: D-0612
title: A host API method the dashboard calls touches writes_methods.go to register it, guardreads_test.go to list it as a read or a write, and flaiover's agent.ts REQUIRED_METHODS, which hostapi's contract_test.go checks against flai's table
kind: decision
sources: [S-0757, T-3181]
topics: [flai, flaiover, dashboard, planning]
paths: [flai/internal/hostapi/writes_methods.go, flai/internal/hostapi/guardreads_test.go, flai/internal/hostapi/contract_test.go, flaiover/src/lib/server/agent.ts]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: []
---

# D-0612 A host API method the dashboard calls touches writes_methods.go to register it, guardreads_test.go to list it as a read or a write, and flaiover's agent.ts REQUIRED_METHODS, which hostapi's contract_test.go checks against flai's table

S-0757 added the host API methods `images.get`, `images.log`, and `images.ensure` for the dashboard's Agent images page (T-3181). The planner's touches named only the new `images.go` and its test, and the task had to widen. Its narrative, `wip/archive/agents/S-0757.md`, records under Decisions:

> T-3181 widened to writes_actions.go (Host.Images), writes_methods.go (registers imageSpecs), serve.go (wires the runner), guardreads_test.go (lists the new methods), and flaiover's agent.ts REQUIRED_METHODS, which is the method list contract_test.go reads. Methods: images.get and images.log, reads; images.ensure, a write under the agent host action.

So a planner or a story's agent adding a host API method the dashboard calls puts these files in the task's touches from the start:

- `flai/internal/hostapi/writes_methods.go`, or the file whose spec list registers it.
- `flai/internal/hostapi/guardreads_test.go`, which lists every method that is not in `guard.HostAPIReads` as a write, with its reason.
- `flaiover/src/lib/server/agent.ts`, whose `REQUIRED_METHODS` `flai/internal/hostapi/contract_test.go` reads and checks against flai's table.

A method with state in flai serve's process also touches `writes_actions.go`, for its `Host` field, and `flai/internal/serve/serve.go`, which wires it.
