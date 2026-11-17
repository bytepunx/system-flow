---
id: D-0612
title: A host API method the dashboard calls touches writes_methods.go to register it, guardreads_test.go to list it as a read or a write, and flaiover's agent.ts REQUIRED_METHODS, which hostapi's contract_test.go checks against flai's table
kind: decision
sources: [S-0757, T-3181, S-0785, T-3354, S-0804, T-3430]
topics: [flai, flaiover, dashboard, planning]
paths: [flai/internal/hostapi/writes_methods.go, flai/internal/hostapi/guardreads_test.go, flai/internal/hostapi/contract_test.go, flaiover/src/lib/server/agent.ts, flai/internal/hostapi/writes_actions.go, flai/internal/serve/serve.go, flai/internal/guard/guard_tables.go, flai/internal/hostapi/writes_test.go]
date: 2026-11-13
status: active
absorbed_by: ""
asked_again: []
cited: [S-0785, T-3354]
---

# D-0612 A host API method the dashboard calls touches writes_methods.go to register it, guardreads_test.go to list it as a read or a write, and flaiover's agent.ts REQUIRED_METHODS, which hostapi's contract_test.go checks against flai's table

S-0757 added the host API methods `images.get`, `images.log`, and `images.ensure` for the dashboard's Agent images page (T-3181). The planner's touches named only the new `images.go` and its test, and the task had to widen. Its narrative, `wip/archive/agents/S-0757.md`, records under Decisions:

> T-3181 widened to writes_actions.go (Host.Images), writes_methods.go (registers imageSpecs), serve.go (wires the runner), guardreads_test.go (lists the new methods), and flaiover's agent.ts REQUIRED_METHODS, which is the method list contract_test.go reads. Methods: images.get and images.log, reads; images.ensure, a write under the agent host action.

So a planner or a story's agent adding a host API method the dashboard calls puts these files in the task's touches from the start:

- `flai/internal/hostapi/writes_methods.go`, or the file whose spec list registers it.
- `flai/internal/hostapi/guardreads_test.go`, which lists every method that is not in `guard.HostAPIReads` as a write, with its reason.
- `flaiover/src/lib/server/agent.ts`, whose `REQUIRED_METHODS` `flai/internal/hostapi/contract_test.go` reads and checks against flai's table.

A method with state in flai serve's process also touches `writes_actions.go`, for its `Host` field, and `flai/internal/serve/serve.go`, which wires it.

S-0785 confirmed it when it added `log.read` (T-3354, [ADR-0300](../adrs/0300-the-dashboard-reads-a-log-flai-wrote-through-log-read-confined-to-the-project-s.md)). Its plan, `wip/archive/agents/S-0785.md`, named these files from the start, citing this entry:

> Plan: layer 1 holds T-3354 (Go host method log.read, plus REQUIRED_METHODS in flaiover's agent.ts, which contract_test.go checks, D-0612) …

T-3354's touches also took `writes_actions.go` and `serve.go`, for the `Host.ServeDir` that flai serve sets so `log.read` can read its state folder.

S-0804 met it again when it added the host read `resolve.status` (T-3430, [ADR-0313](../adrs/0313-rebase-and-resolve-runs-as-a-flai-serve-job-that-rebases-a-review-story-s.md)). Its narrative, `wip/archive/agents/S-0804.md`, records under Decisions:

> T-3430 widened to the method lists a new host read must join: guard_tables.go's HostAPIReads, hostapi writes_test.go's good and refused tables, guardreads_test.go's comment, and flaiover's REQUIRED_METHODS in agent.ts, which the contract test checks; made by the story's agent after the sub-agent named them.

So a new host read also touches:

- `flai/internal/guard/guard_tables.go`, whose `HostAPIReads` lets every agent run it through `flai hostapi` (D-0538).
- `flai/internal/hostapi/writes_test.go`, whose good and refused tables list it.
