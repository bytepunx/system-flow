---
id: S-0244
type: story
nature: remediation
title: "Two in-progress stories whose claims grew to overlap each wait for the other: close-out fails flai check --strict on wip.overlap, and the wait costs a model turn every five minutes"
status: in-progress
owner: alex
created: 2026-10-03T03:31:01Z
updated: 2026-10-05T05:03:51Z
transitions:
  - to: ready
    at: 2026-10-04T21:44:16Z
    by: alex
  - to: in-progress
    at: 2026-10-05T04:41:19Z
    by: agent-S-0244
tags: [flai, template]
touches: [flai/internal/workitem/hold.go, flai/internal/itemedit, flai/internal/itemnew, flai/cmd/items.go, flai/cmd/edit.go, flai/cmd/touches.go, flai/internal/mcpserver, flai/cmd/mcp.go, flai/cmd/mcp_http.go, flai/internal/harness, design/conventions/delegation.md, design/conventions/work-management.md, template/root/design/conventions/delegation.md, template/root/design/conventions/work-management.md, template/CHANGELOG.md, design/system/agent-narrative.md, design/system/workflow.md, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md, design/issues/I-0059-two-in-progress-stories-whose-claims-grew-to-overlap-each-wait-for-the-other-close-out-fails-flai-check-strict-on-wip-overlap-and-the-wait-costs-a-model-turn-every-five-minutes.md, design/issues/summary.md, flai/cmd/touches_overlap_test.go, template/template.yaml]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 690
  models:
    - model: claude-opus-5-5
      input: 218
      output: 86627
      cache_read: 9432793
      cache_write: 347265
      cost: 5.5645
cost_of_delay:
  inputs:
    penalty_per_week: 625
    by: planner-S-0244
    at: 2026-10-05T03:49:11Z
  value: 625
  by: planner-S-0244
  at: 2026-10-05T03:49:16Z
forecast:
  duration: 1h30m
  delivery: 2026-10-05T06:25:00Z
  basis: "Its own forecast of 1h30m; 1st in the pull order with an in-progress limit of 3, behind S-0257."
  by: flai
  at: 2026-10-05T04:38:16Z
---
# S-0244 Two in-progress stories whose claims grew to overlap each wait for the other: close-out fails flai check --strict on wip.overlap, and the wait costs a model turn every five minutes

## Goal

This story remediates [I-0059](../../../design/issues/I-0059-two-in-progress-stories-whose-claims-grew-to-overlap-each-wait-for-the-other-close-out-fails-flai-check-strict-on-wip-overlap-and-the-wait-costs-a-model-turn-every-five-minutes.md), "Two in-progress stories whose claims grew to overlap each wait for the other: close-out fails flai check --strict on wip.overlap, and the wait costs a model turn every five minutes". The issue recommends this solution:

Possible remediations, cheapest first:

1. **The close-out does not fail on another story's claim.** `wip.overlap` between two in-progress stories is the pull-time hold's business and advisory afterwards (ADR-0019, ADR-0046); `scripts/close-out.sh` (here and in the template) should not fail on it: pass an allow list to `flai check --strict` (`--allow wip.overlap`), or give `flai check` a `--story S-nnnn` scope that reports only what the story can clear. The delegation and work-management conventions should say that an overlap with another open story is not the story's to clear: note it in the narrative and go to review.
2. **Report a claim that grows into another story's.** When `flai task new`, `flai edit --touches`, or the MCP writes give an open story's task touches that overlap another in-progress story's claim, say so in the command's output and send both stories an inbox event (`overlapped`, as acceptance already does, S-0132), so the agents coordinate while the work is small instead of discovering it at close-out.
3. **Let a long poll be long.** `wait_for_events` and `wait_for_work` should honour the requested `timeout_seconds` up to what the transport allows (the agents asked for 900 and 1800 s and got 300), so an idle agent costs one turn per answer, not one per five minutes. Measure: at 150k cached tokens a turn, a five-minute tick costs about 0.05 USD; an hour of waiting about 0.60 USD per agent plus the turns it spends deciding to wait again.
4. **Do not wait for sub-agents with flai.** The delegation convention and `harness.Prompt` should tell the story's agent to wait for a background sub-agent through the harness's own completion notice, not by polling `wait_for_events`, which cannot see a sub-agent end.
5. **Detect the mutual wait.** `flai serve` or the designer's inbox could report two in-progress stories each blocked only on the other's claim (both agents idle in waits, both close-outs failing on `wip.overlap` against each other) and open a thread to the operator naming the pair and what clears it, rather than leaving two agents polling for an hour.

Related: S-0197 (stream sync and conflict reporting), the agent-waiting chart in E-0016 (S-0215) would have made this visible; ADR-0046 for the hold.

## Acceptance criteria
- [ ] The cause I-0059 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0059 is closed with `flai issue close I-0059 --reason` saying what fixed it

## Tasks
- T-0867 A write that grows an open story's claim into an in-progress story's says so and tells both stories as an overlapped change
- T-0868 wait_for_events and wait_for_work honour the requested timeout up to a 30-minute cap instead of five minutes
- T-0869 The story agent's prompt and the conventions say to wait for sub-agents through the harness's notice, and what to do on a grown-claim overlapped change
- T-0870 The design and the user docs describe the grown-claim notice and the longer wait, and I-0059 is closed

## Notes

### Planning

Planned by planner-S-0244 on 2026-10-05, plan thread TH-0125; the cost of delay question is TH-0121.

**State of I-0059's remediations.** Remediation 1 is already done: S-0249 scoped the close-out's `flai check --strict` to the story, and a `wip.overlap` is now a note recorded in an issue (`scripts/close-out.sh`, work-management.md). The mutual wait at close-out cannot happen any more. What is left is remediations 2, 3, and 4. Remediation 5, detecting the mutual wait, is not planned: once the close-out no longer stops on `wip.overlap`, two stories can no longer each wait for the other.

**Touches.** The story declared none; every touch is predicted.

- `flai/internal/workitem/hold.go`, `flai/internal/itemedit`, `flai/internal/itemnew`, `flai/cmd/items.go`, `flai/cmd/edit.go`, `flai/cmd/touches.go`, and `flai/internal/mcpserver`. From the layout: the claim and overlap live in `hold.go`, and `itemedit.RecordOverlap` and the `overlapped` reporting are in `mcpserver/cursor.go`. The touches writers are the three commands and `mcpserver/items_write.go`. The dashboard writes through the CLI (`hostapi/writes.go`), so it is not touched. Co-change confirms `mcpserver` and `cursor.go`.
- `flai/cmd/mcp.go` and `flai/cmd/mcp_http.go`. From the layout: the server options, and the comment that holds the five-minute wait. Co-change lists `mcp.go`.
- `flai/internal/harness`. From the design (the issue names `harness.Prompt`). Co-change confirms `harness_test.go`.
- `design/conventions/delegation.md` and `design/conventions/work-management.md`, and their copies under `template/root/design/conventions`, with `template/CHANGELOG.md`. From the design (the issue names both conventions). Co-change lists the template copy of delegation.md, `work-management.md`, and the changelog.
- `design/system/agent-narrative.md`, `workflow.md`, `flai-cli.md`, `docs/users/flai.md`, and `flai-reference.md`. From the design (these describe `overlapped` and `wait_for_events`). Co-change lists `flai.md`, `flai-cli.md`, and `flai-reference.md`.
- I-0059 and `design/issues/summary.md`, from the second criterion.
- Not taken: `flai/internal/serve/agents.go`, `cmd/serve_actions.go`, and `docs/operators/settings.md` came up in co-change only through the harness.

**Tasks and layers.**

1. T-0867 (grown-claim notice) and T-0868 (long-poll cap) run together.
2. T-0869 (prompt and conventions) waits for both.
3. T-0870 (design, docs, and I-0059 closed) waits for all three.

**Forecast.** `flai forecast` gave 24m, from 60 s per unit over 8 done large remediation stories, times size 24. Raised to 1h30m. Four tasks, and T-0867 is a new notice across five writers and the inbox cursor, with tests. S-0249 and S-0266, which were smaller, took about 30 to 50 minutes each. Delivery is 06:00Z: the story is held on overlap behind S-0253 and S-0260, both in progress since 03:10Z to 03:14Z, and recent stories have been accepted about 30 minutes after review.

**Cost of delay.** The value is 625.00 USD a week, from the operator's `penalty_per_week` of 625, given in TH-0121 ("go with the alternative"), and `flai cod`'s value stands as computed. The input charges for the model turns the waits burned. I-0059's first instance cost 12.68 USD (S-0198) and 3.14 USD (S-0242), about 15.82 USD an incident. There were 4 incidents in about 17 hours on 2026-10-03, about 39.5 a week, which comes to about 625 USD a week. S-0249 has since removed the close-out stop, so the rate is now likely lower; the figure is the operator's choice.
