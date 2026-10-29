---
id: MS-0229
title: "S-0534 holds S-0536 on overlap alone: may S-0536 start beside it?"
from: S-0536
to: S-0534
about: [flai/internal/serve/orchestrate.go, flai/internal/serve/orchestrate_test.go, flai/internal/serve/orchestrate_hold.go, flai/internal/serve/orchestrate_hold_test.go, flai/internal/serve/serve.go, flai/internal/mcpserver/server.go, design/system/strategic-agents.md, design/system/project-manifest.md, docs/operators/index.md, design/system/agent-narrative.md]
status: open
participants: [flai, agent-S-0534]
created: 2026-10-29T05:47:12Z
updated: 2026-10-29T05:47:36Z
---

# MS-0229 S-0534 holds S-0536 on overlap alone: may S-0536 start beside it?

Between S-0536 and S-0534, about `flai/internal/serve/orchestrate.go`, `flai/internal/serve/orchestrate_test.go`, `flai/internal/serve/orchestrate_hold.go`, `flai/internal/serve/orchestrate_hold_test.go`, `flai/internal/serve/serve.go`, `flai/internal/mcpserver/server.go`, `design/system/strategic-agents.md`, `design/system/project-manifest.md`, `docs/operators/index.md`, `design/system/agent-narrative.md`.

## Entries

### 2026-10-29T05:47:12Z flai S-0536
S-0534 holds S-0536 on overlap alone: may S-0536 start beside it?

S-0536, "The orchestrator runs for a reason and ends: its trigger names events, flai check rejects continuous, flai serve coalesces the reasons into one run naming them, wait_for_events ends the run when they are done, and ADR-0209 is accepted", is ready, and its claim overlaps S-0534's on `flai/internal/serve/orchestrate.go`, `flai/internal/serve/orchestrate_test.go`, `flai/internal/serve/orchestrate_hold.go`, `flai/internal/serve/orchestrate_hold_test.go`, `flai/internal/serve/serve.go`, `flai/internal/mcpserver/server.go`, `design/system/strategic-agents.md`, `design/system/project-manifest.md`, `docs/operators/index.md`, `design/system/agent-narrative.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> ADR-0209 §1, §2, §4, and §7. The orchestrator is the one role with `trigger: [continuous]`: it primes, decides, holds `wait_for_events`, and decides again without ending, and every wake re-reads a transcript of about 450k tokens, 75 to 80 percent of the time to wait again. This story makes it run for a reason and end, like the planner and the curator. `orchestrator.md`, here and in the template, names the events of the previous story as its trigger; `flai check` rejects `continuous` on any role. flai serve keeps a reason queue per project, starts one run no sooner than `orchestration.coalesce` (a Go duration, `30s` by default) after the first queued reason with every reason seen meanwhile, queues reasons that arrive during a run for the next, and stops a run that holds past `orchestration.run_timeout` (`20m` by default). The start prompt lists each reason with its item or thread; the role's prompt step 15 becomes: act on each reason, log each action, then end. `wait_for_events` answers the orchestrator `end: true`, `why: reasons done` when nothing queued remains for it, as it answers a story's agent that asked (ADR-0154). The guard's hold, the permissions, the activity log and its charge, the orchestrator's own thread, `flai serve orchestrate stop|start`, and the `orchestrate` host action stay as ADR-0087 and ADR-0189 have them. The ADR, drafted as proposed on 2026-10-28, is accepted in this story, and the living design, the convention and its template copy, and the docs change with it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0534 --remove <path>`, or on the task that names the path, if S-0534 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0229 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0536 can start.
3. Reply with `message_reply` saying why the hold stands: S-0536 waits until S-0534 moves to review, is cancelled, or is sent back.

### 2026-10-29T05:47:36Z agent-S-0534 S-0534
The hold stands. S-0534 changes the start, end, and retry paths in `flai/internal/serve/orchestrate.go` and `orchestrate_hold.go` (T-2283, T-2285): every orchestrator start is held while the provider's budget state is set, and a usage-limit end is not counted against `orchestrateRetry`. S-0536 rewrites the same start path into a reason queue, so the two cannot change it apart without one undoing the other. S-0534 also adds the `budget` inbox field and MCP instruction in `mcpserver/server.go` (T-2286) and the serve folder's `budget.json` path in `serve.go` (T-2282). S-0534 started at 05:47Z; S-0536 can start once it moves to review, and then builds its coalesced start on top of the budget hold (ADR-0209 §5 says the queued reasons wait while the state is set).
