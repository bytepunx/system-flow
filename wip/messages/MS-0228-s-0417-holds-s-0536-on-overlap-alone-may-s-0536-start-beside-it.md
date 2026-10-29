---
id: MS-0228
title: "S-0417 holds S-0536 on overlap alone: may S-0536 start beside it?"
from: S-0536
to: S-0417
about: [flai/internal/serve/serve.go, flai/internal/harness/harness.go, flai/internal/harness/harness_test.go, flai/internal/manifest/manifest.go, flai/internal/manifest/manifest_test.go, flai/internal/manifest/settings.go, flai/internal/manifest/settings_test.go, design/system/strategic-agents.md, design/system/project-manifest.md, docs/operators/settings.md, docs/operators/index.md, design/system/agent-narrative.md]
status: open
participants: [flai, agent-S-0417]
created: 2026-10-29T04:22:22Z
updated: 2026-10-29T04:54:03Z
---

# MS-0228 S-0417 holds S-0536 on overlap alone: may S-0536 start beside it?

Between S-0536 and S-0417, about `flai/internal/serve/serve.go`, `flai/internal/harness/harness.go`, `flai/internal/harness/harness_test.go`, `flai/internal/manifest/manifest.go`, `flai/internal/manifest/manifest_test.go`, `flai/internal/manifest/settings.go`, `flai/internal/manifest/settings_test.go`, `design/system/strategic-agents.md`, `design/system/project-manifest.md`, `docs/operators/settings.md`, `docs/operators/index.md`, `design/system/agent-narrative.md`.

## Entries

### 2026-10-29T04:22:22Z flai S-0536
S-0417 holds S-0536 on overlap alone: may S-0536 start beside it?

S-0536, "The orchestrator runs for a reason and ends: its trigger names events, flai check rejects continuous, flai serve coalesces the reasons into one run naming them, wait_for_events ends the run when they are done, and ADR-0209 is accepted", is ready, and its claim overlaps S-0417's on `flai/internal/serve/serve.go`, `flai/internal/harness/harness.go`, `flai/internal/harness/harness_test.go`, `flai/internal/manifest/manifest.go`, `flai/internal/manifest/manifest_test.go`, `flai/internal/manifest/settings.go`, `flai/internal/manifest/settings_test.go`, `design/system/strategic-agents.md`, `design/system/project-manifest.md`, `docs/operators/settings.md`, `docs/operators/index.md`, `design/system/agent-narrative.md`. It has no agent until it starts, so flai asks for it (ADR-0134).

Its goal:

> ADR-0209 §1, §2, §4, and §7. The orchestrator is the one role with `trigger: [continuous]`: it primes, decides, holds `wait_for_events`, and decides again without ending, and every wake re-reads a transcript of about 450k tokens, 75 to 80 percent of the time to wait again. This story makes it run for a reason and end, like the planner and the curator. `orchestrator.md`, here and in the template, names the events of the previous story as its trigger; `flai check` rejects `continuous` on any role. flai serve keeps a reason queue per project, starts one run no sooner than `orchestration.coalesce` (a Go duration, `30s` by default) after the first queued reason with every reason seen meanwhile, queues reasons that arrive during a run for the next, and stops a run that holds past `orchestration.run_timeout` (`20m` by default). The start prompt lists each reason with its item or thread; the role's prompt step 15 becomes: act on each reason, log each action, then end. `wait_for_events` answers the orchestrator `end: true`, `why: reasons done` when nothing queued remains for it, as it answers a story's agent that asked (ADR-0154). The guard's hold, the permissions, the activity log and its charge, the orchestrator's own thread, `flai serve orchestrate stop|start`, and the `orchestrate` host action stay as ADR-0087 and ADR-0189 have them. The ADR, drafted as proposed on 2026-10-28, is accepted in this story, and the living design, the convention and its template copy, and the docs change with it.

Answer in one of three ways:

1. Narrow your touches with `flai touches S-0417 --remove <path>`, or on the task that names the path, if S-0417 will not change it: the hold clears when the claims no longer overlap.
2. Share the paths with `flai message share MS-0228 --paths <path> "<split>"`, or the MCP tool `message_share`, saying who changes what, if the two stories can change them apart: the overlap on them no longer holds, and S-0536 can start.
3. Reply with `message_reply` saying why the hold stands: S-0536 waits until S-0417 moves to review, is cancelled, or is sent back.

### 2026-10-29T04:54:03Z agent-S-0417 S-0417 (settled)
The hold stands. S-0417 has changed `flai/internal/serve/serve.go`, `harness/harness.go`, and `manifest/manifest.go` and `settings.go`: it adds defined roles' runs, their harness requests, and `roles.<name>.agent`. S-0536 changes the same files. S-0417 is in its last task, the docs, and moves to review within the hour, which clears the hold.
