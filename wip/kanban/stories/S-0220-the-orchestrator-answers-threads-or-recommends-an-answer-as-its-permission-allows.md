---
id: S-0220
type: story
nature: feature
title: The orchestrator answers threads, or recommends an answer, as its permission allows
status: ready
parent: E-0016
owner: arobson
created: 2026-10-02T11:54:16Z
updated: 2026-10-05T06:13:40Z
transitions:
  - to: ready
    at: 2026-10-05T04:41:20Z
    by: alex
tags: [flai]
touches: [flai/internal/harness, ".claude/agents/orchestrator.md", flai/internal/threads, design/system/strategic-agents.md, flai/internal/mcpserver, flai/internal/hostapi, flai/internal/metrics, design/system/metrics.md, design/adrs, flaiover/src/routes/inbox, flaiover/src/routes/threads, flaiover/src/routes/api/inbox, design/system/flaiover-dashboard.md, docs/users/flaiover.md, docs/users/flai.md, flai/cmd/thread.go, flai/cmd/thread_test.go, flai/internal/guard, flai/internal/serve/agents.go, flai/internal/serve/agents_test.go, design/system/flai-cli.md, docs/users/flai-reference.md, flaiover/src/routes/api/threads, flaiover/src/lib/components/Threads.svelte, flaiover/src/lib/components/Threads.svelte.test.ts, flaiover/src/lib/components/InboxView.svelte, flaiover/src/lib/components/Inbox.svelte.test.ts, flaiover/src/lib/server/inbox.ts, flaiover/src/lib/server/inbox.test.ts, template/root/.claude/agents/orchestrator.md, design/conventions/strategic-agents.md, template/root/design/conventions/strategic-agents.md, template/CHANGELOG.md]
after: [S-0218]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
cost_of_delay:
  value: 76.01
  by: planner-S-0220
  at: 2026-10-05T04:49:47Z
forecast:
  duration: 1h15m
  delivery: 2026-10-05T13:35:00Z
  basis: "Its own forecast of 1h15m; 3rd in the pull order with an in-progress limit of 3, behind S-0217, S-0218 and S-0219."
  by: flai
  at: 2026-10-05T06:13:40Z
---
# S-0220 The orchestrator answers threads, or recommends an answer, as its permission allows

## Goal

Agents wait on threads for the operator. With `answer_threads: recommend`, the orchestrator replies to each thread awaiting the operator with its recommended answer, marked as a recommendation, for the operator to confirm; with `autonomous`, it answers and the agent goes on.

## Acceptance criteria
- [ ] In `recommend`, the orchestrator posts a reply tagged as a recommendation that does not set the thread to `answered`, citing the design, ADR, or convention it drew on; the inbox shows the operator a thread with a recommendation to confirm with one action (confirm makes it the answer)
- [ ] In `autonomous`, it answers the thread (setting `answered`) when it can cite a source, and escalates to the operator with a recommendation when it cannot or when the question names the operator's judgement (a decision, a scope change, money)
- [ ] It never resolves a thread it did not open, and never answers a thread opened by itself
- [ ] Each answer is logged with the source cited; the metrics count thread waits ended by the orchestrator separately
- [ ] Tests cover a recommendation, an autonomous answer, and an escalation

## Tasks
- T-0894 A thread entry can be a recommendation with a source, which leaves the thread awaiting the operator until they confirm it
- T-0902 flai thread reply, the MCP tool thread_reply, and the host API post a recommendation with its source, and flai thread confirm makes it the answer
- T-0906 flai stats counts the thread waits the orchestrator ended apart from the operator's, and a recommendation ends no wait
- T-0908 A pending recommendation is not an answer to the agent that asked, and the designer's inbox lists it to confirm
- T-0910 flai guard holds the orchestrator's thread calls to answer_threads: recommend only, an answer only with a source, never resolving another's thread or answering its own
- T-0912 The orchestrator's prompt and definition tell it to recommend, answer, or escalate each thread awaiting the operator, citing its source
- T-0913 The dashboard's inbox and threads show a pending recommendation with its source and confirm it with one action
- T-0914 An end-to-end test drives the orchestrator's thread calls through the MCP server and the guard: a recommendation, an autonomous answer, and an escalation

## Notes

### Planning

Planned by planner-S-0220 on 2026-10-05. The plan thread is TH-0136, with the layers and nine assumptions to confirm.

Touches, by source:

- Declared, all kept: `flai/internal/harness`, `.claude/agents/orchestrator.md`, `flai/internal/threads`, `design/system/strategic-agents.md`, `flai/internal/mcpserver`, `flai/internal/hostapi`, `flai/internal/metrics`, `design/system/metrics.md`, `design/adrs`, `flaiover/src/routes/inbox`, `flaiover/src/routes/threads`, `flaiover/src/routes/api/inbox`, `design/system/flaiover-dashboard.md`, `docs/users/flaiover.md`, `docs/users/flai.md`.
- Co-change (`flai touches suggest`): `design/system/flai-cli.md` (40% of the 412 commits, and the thread format and commands are described there) and `docs/users/flai-reference.md` (15%, the `flai thread` reference). The other suggestions, such as `docs/operators/index.md`, `flai/cmd/serve_actions.go`, and `flai/internal/check/check.go`, follow S-0218's host action and checks, not threads, and are left out.
- Design: `design/conventions/strategic-agents.md`, `template/root/design/conventions/strategic-agents.md`, and `template/CHANGELOG.md`, because the convention's `## As the orchestrator` lists answering threads and S-0210 set the precedent of changing the convention with its template copy.
- Layout: `flai/cmd/thread.go` and `flai/cmd/thread_test.go` (`flai thread reply` and the new `confirm`); `flai/internal/guard` (the orchestrator's thread rules); `flai/internal/serve/agents.go` and its test (`answered` starts an agent again on any later entry by another author, so it would start one on a recommendation); `flaiover/src/routes/api/threads` (the confirm route sits beside `reply` and `resolve`); `flaiover/src/lib/components/Threads.svelte`, `InboxView.svelte`, `flaiover/src/lib/server/inbox.ts`, and their tests (where the inbox and the threads draw); `template/root/.claude/agents/orchestrator.md` (the template copy of the definition, as `planner.md` has).

Forecast: 1h15m, adjusted. flai's history rate gives 45m on the 15 declared touches and 1h24m on the 33 predicted, which list test files apart and so overstate the size. Six done stories that changed Go, the dashboard, and an ADR together took 53 to 72 minutes of agent time (S-0201, S-0204, S-0205, S-0208, S-0210, S-0225). This one has eight tasks in four layers and two ADRs, so it sits at the top of that range. The delivery is flai's play-out (6th in the pull order, after S-0218) with that duration.

Cost of delay: 76.01 USD a week, worked out by `flai cod`, which stands. S-0220 has no inputs of its own, so it takes its share of E-0016's 1500 USD a week (the operator's 10h lost a cycle at 150 USD an hour): 1h15m of the 24h40m forecast over the epic's 17 open stories without inputs. It was 97.56 from the epic plan's 2h forecast; it fell with the shorter forecast.
