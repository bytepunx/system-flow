---
id: TH-1080
title: "S-0714 plan: four tasks in four layers, 45m, 25 USD a week"
anchor:
  path: wip/kanban/stories/S-0714-testeverylistedflaiovercomponentcompilesintheharness-fails-a-close-out-with-the.md
  item: S-0714
status: resolved
participants: [planner-S-0714, orchestrator]
created: 2026-11-11T09:06:40Z
updated: 2026-11-11T09:09:23Z
---

# TH-1080 S-0714 plan: four tasks in four layers, 45m, 25 USD a week

On wip/kanban/stories/S-0714-testeverylistedflaiovercomponentcompilesintheharness-fails-a-close-out-with-the.md.

## Entries

### 2026-11-11T09:06:40Z planner-S-0714
No earlier plan or thread existed on S-0714 or I-0342, so this is the first plan.

Proposed cause: the 404 body `404 page not found` is Go's net/http default, not Vite's. A Go server answered on the port where vite had answered the harness page with 200 a moment before. `freePort` takes vite's port from a listen on 127.0.0.1:0, which draws from the kernel's ephemeral range (32768 to 60999 here). Every httptest server and outgoing connection of the other close-outs running at the same time draws from that range too. Once vite stopped listening, by exiting or restarting, one of them could be given the port.

Tasks, one per layer, each waiting on the one before:

1. T-3114: freePort picks the port below the host's ephemeral range (devserver.go, devserver_test.go).
2. T-3115: the harness plugin marks every answer with a per-start token, and StartDevServer's wait takes only a 200 that carries it. A test reproduces the cause with a Go server on the picked port. It waits on T-3114 because both change devserver.go and devserver_test.go.
3. T-3116: TestEveryListedFlaioverComponentCompilesInTheHarness checks the token on each answer. When the token is missing, it says another server answered and gives vite's state, instead of reporting a compile failure. It reads T-3115's token and header.
4. T-3117: design/system/component-library.md describes the port's range and the token, and I-0342 is closed. It waits until the fix and its test are in.

Figures:

- Forecast: 45m, up from flai's 27m. The fix spans Go and the JavaScript plugin, and its proof needs real vite runs in the integration tier.
- Cost of delay: 25 USD a week, from flai's own 10m time_lost_per_cycle input. I left that input unchanged.

Assumptions:

- The cause is inferred from one instance whose vite last lines were not recorded. The story's agent should confirm it before building, as the goal asks. If vite never stopped listening, T-3115's token still turns the failure into a clear report, but T-3114 may not be the fix.
- No platform keeps a port below its ephemeral range from being free. Off Linux, the IANA range 49152 to 65535 is taken as the ephemeral one.
- No user-facing behaviour changes, so docs/users is not touched.

### 2026-11-11T09:09:09Z orchestrator
Plan approved. T-3114, T-3115, T-3116, and T-3117 run one per layer, each waiting on the one before. Every task's touches fall within the story's. The 45m forecast and the 25 USD a week value, from flai's own 10m input, fit the story. Nothing here asks a question.

The assumption stands that the story's agent confirms the cause before building, as the goal asks.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-11T09:09:23Z orchestrator
Resolved: Plan approved; no question asked.
