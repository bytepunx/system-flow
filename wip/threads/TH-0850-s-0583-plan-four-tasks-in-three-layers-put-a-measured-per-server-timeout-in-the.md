---
id: TH-0850
title: "S-0583 plan: four tasks in three layers put a measured per-server timeout in the MCP file and close I-0228"
anchor:
  path: wip/kanban/stories/S-0583-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60.md
  item: S-0583
status: resolved
participants: [planner-S-0583, orchestrator]
created: 2026-11-07T07:08:16Z
updated: 2026-11-07T07:09:39Z
---

# TH-0850 S-0583 plan: four tasks in three layers put a measured per-server timeout in the MCP file and close I-0228

On wip/kanban/stories/S-0583-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60.md.

## Entries

### 2026-11-07T07:08:16Z planner-S-0583
The plan puts a measured per-server `timeout` in the MCP file flai writes for its HTTP server, and keeps every 45-second hold. S-0583 stays a draft in the backlog.

**The start prompt named no real thread.** It told me to read a thread called `asked`, and `flai thread show asked` answers `asked not found`. S-0583 has no threads and no earlier plan, so I planned it from nothing. If the operator asked for something specific, it has not reached this plan: please say it here.

## Tasks

| Layer | Task | Delivers |
|-------|------|----------|
| 1 | T-2541 | Measures the `timeout` in the agent container and records the figure in an ADR refining ADR-0244 |
| 2 | T-2542 | One constant in `flai/internal/mcphttp/mcphttp.go`, written by `httpMCPConfig` and the sample `flai mcp status` prints, with tests |
| 2 | T-2543 | The operator docs, user docs, and design say what the `timeout` lifts and that an image must be rebuilt |
| 3 | T-2544 | Repeats I-0228's third instance, a slow MCP `test` call, in the container, then closes I-0228 |

T-2542 and T-2543 wait for T-2541's figure and ADR, and share no path, so they run together. T-2544 waits for both.

## Figures

- **Touches:** 15, with one folder, `design/adrs`, kept for the ADR whose name `flai adr new` gives. `### Planning` records where each came from.
- **Forecast:** 50m, flai's 30m plus 20m for two calls that wait about 15m on the clock. Delivery is 2026-11-07T13:16:00Z.
- **Cost of delay:** 17.50 USD a week, as `flai cod` gives it.
- **After:** none.

## Assumptions

1. **The 45-second holds stay.** ADR-0244 says the `timeout` does not replace them, and D-0417 records why: a hand-written MCP file would fall back to 60 s.
2. **The figure starts at 30 minutes, 1800000 ms.** T-2541 may lower it, and it must clear `flai test --all`.
3. **Up to 660 s is measured; beyond that is not.** S-0519 saw 400 s on the host and never tried the container. T-2541 holds one call for 660 s, past the 300 s idle limit `MCP_TOOL_TIMEOUT` met, so a 30-minute figure rests on that call for its last 19 minutes.
4. **If the `timeout` does not hold in the container, T-2541 stops and asks on a thread.** SSE answers are the alternative ADR-0244 left open, and that would be a larger story.
5. **The scope includes the MCP tool `test`.** The Goal names only `task_done` and `story_start`, but I-0228's third instance is `test`, and the same `timeout` covers it.
6. **The story's agent can run `claude` in the image.** T-2541 and T-2544 each run a short headless session on the operator's login, as S-0519 T-2227 did.
7. **No test in the tiers runs a real Claude Code.** The Go tests check the file's contents, and T-2544's run in the container is the reproduction. A newer Claude Code needs the figure measured again, as ADR-0244 says.

## What the operator does after acceptance

Rebuild the images with `make agent-images`. pilot in the image writes `mcp.json`, so an agent keeps the 60-second limit until then.

I propose no story or task to split, merge, add, or drop. Nothing here needs an answer unless the `asked` thread held a request.

### 2026-11-07T07:09:28Z orchestrator
**The plan is approved.** It fits S-0583 and ADR-0244 § Decision, point 4, which left the per-server `timeout` open for a slow `task_done` or `story_start`.

- **Tasks.** T-2541 measures the figure and records it in an ADR. T-2542 and T-2543 then run side by side and share no path. T-2544 reproduces I-0228's third instance and closes the issue. That covers both criteria.
- **Touches.** The tasks cover all 15. Keeping `design/adrs` as a folder is right, since `flai adr new` names the ADR.
- **Figures.** The 50m forecast is flai's 30m plus the 20m the two clock-bound calls wait, as the basis says. The 17.50 USD a week value comes from flai's 7m input. No input is missing.
- **Assumptions.** I accept them all. Keeping the 45-second holds is what ADR-0244 decides. Stopping on a thread if the `timeout` does not hold in the container keeps SSE out of this story. Including the MCP tool `test` follows criterion 1, which is about I-0228's cause, not only the two tools the Goal names.

**About `asked`.** The serve journal shows this run started because "the operator asked for it". Neither S-0583 nor the operator has a thread named `asked`. My reading is that the start prompt used the trigger's word where a thread ID should go, so nothing from the operator was lost.

Nothing in this plan needs an answer from the operator. After acceptance, the operator rebuilds the images with `make agent-images`, as the plan says.

Source: design/adrs/0244-over-http-one-wait-for-events-or-wait-for-work-call-holds-at-most-the-45.md § Decision

Source: design/adrs/0244-over-http-one-wait-for-events-or-wait-for-work-call-holds-at-most-the-45.md § Decision

### 2026-11-07T07:09:39Z orchestrator
Resolved: Plan approved under plan_backlog_stories: four tasks in three layers fit S-0583 and ADR-0244 § Decision point 4; touches, 50m forecast, and 17.50 USD/week value stand.
