---
id: TH-0125
title: "S-0244 plan: four tasks in three layers; I-0059's remediations 2 to 4, 1 already done by S-0249, 5 dropped"
anchor:
  path: wip/kanban/stories/S-0244-two-in-progress-stories-whose-claims-grew-to-overlap-each-wait-for-the-other-close-out-fails-flai-check-strict-on-wip-overlap-and-the-wait-costs-a-model-turn-every-five-minutes.md
  item: S-0244
status: open
participants: [planner-S-0244, alex]
created: 2026-10-05T03:17:16Z
updated: 2026-10-05T03:49:42Z
---

# TH-0125 S-0244 plan: four tasks in three layers; I-0059's remediations 2 to 4, 1 already done by S-0249, 5 dropped

On wip/kanban/stories/S-0244-two-in-progress-stories-whose-claims-grew-to-overlap-each-wait-for-the-other-close-out-fails-flai-check-strict-on-wip-overlap-and-the-wait-costs-a-model-turn-every-five-minutes.md.

## Entries

### 2026-10-05T03:17:16Z planner-S-0244
**Tasks, in three layers:**

1. **Layer 1**, run together, no path in common:
   - **T-0867** A write that grows an open story's claim into an in-progress story's says so and tells both stories as an `overlapped` change. This is remediation 2. It covers `flai task new`, `flai edit --touches`, `flai touches`, and the MCP `item_new` and `item_edit`. It reuses `workitem.Claim`, `PathsOverlap`, and `.flai-cache/overlaps.jsonl`, with a new cause field for an in-progress story.
   - **T-0868** `wait_for_events` and `wait_for_work` honour the requested timeout up to a 30-minute cap instead of five minutes. This is remediation 3. It first finds the claude-code MCP client's tool-call limit.
2. **Layer 2:** **T-0869**, after T-0867 and T-0868. The story agent's prompt (`harness.Prompt`) and `delegation.md` say to wait for background sub-agents through the harness's completion notice, not `wait_for_events`; this is remediation 4. `work-management.md` says what to do on a grown-claim `overlapped` change. Both conventions change here and under `template/root`. If T-0868 found the client's limit too short, this task sets it in the claude-code agent's environment.
3. **Layer 3:** **T-0870**, after all three. The design and the user docs describe the new notice and the wait cap, and I-0059 is closed with each remediation's fix named.

**Assumptions:**

- **Remediation 1 is done.** S-0249 scoped the close-out's check to the story and turned `wip.overlap` into a note. No task repeats it.
- **Remediation 5 (detect the mutual wait) is dropped.** With the close-out no longer stopping on `wip.overlap`, two stories cannot each wait for the other. If you want the detection anyway, it is a story of its own.
- **The grown-claim notice reuses the `overlapped` kind** with `cause` the other in-progress story, rather than a new kind. It reports only the paths a write adds, so an overlap the pull-time hold already allowed is not repeated. It is advisory: the write always stands.
- **30 minutes for the cap** is what the agents in I-0059 asked for (up to 1800 s). T-0868 lowers it if the MCP client cannot hold a call that long.
- I added the tags `flai` and `template`, which flai check requires because the story touches both.

**Figures:**

- **Touches:** 22, all predicted. Where each came from is in the story's Notes under Planning.
- **Forecast:** 1h30m, raised from flai's 24m, with delivery about 06:00Z. It is held behind S-0253 and S-0260 on `design/system/workflow.md` and `flai/internal/itemedit`.
- **Cost of delay:** waits on your answer in TH-0121.

Nothing to split, merge, or drop among the tasks.

### 2026-10-05T03:48:36Z alex
Resolved.

### 2026-10-05T03:49:42Z planner-S-0244
The cost of delay is now set from TH-0121: `penalty_per_week` 625, value **625.00 USD/week**.

flai has since replanned the delivery to 05:33Z, after S-0260's acceptance. It kept the 1h30m duration.

The plan is complete. T-0867 to T-0870 are in the backlog, and `flai check --strict` and the markdown lint are clean on S-0244 and its tasks.
