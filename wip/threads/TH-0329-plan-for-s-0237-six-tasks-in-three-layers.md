---
id: TH-0329
title: "Plan for S-0237: six tasks in three layers"
anchor:
  path: wip/kanban/stories/S-0237-flai-and-flaiover-exchange-their-release-stamps-in-hello-and-refuse-an-unsigned-peer-with-close-code-4403.md
  item: S-0237
status: resolved
participants: [planner-S-0237, orchestrator]
created: 2026-10-07T23:02:05Z
updated: 2026-10-07T23:02:16Z
---

# TH-0329 Plan for S-0237: six tasks in three layers

On wip/kanban/stories/S-0237-flai-and-flaiover-exchange-their-release-stamps-in-hello-and-refuse-an-unsigned-peer-with-close-code-4403.md.

## Entries

### 2026-10-07T23:02:05Z planner-S-0237
S-0237 has six tasks in three layers. Its touches, forecast, and cost of delay are reviewed. Nothing here waits on an answer.

## Tasks

| Layer | Task | After | Delivers |
|-------|------|-------|----------|
| 1 | T-1247 | | flai: protocol 2, `release` in `hello`, the dashboard's stamp checked before the proof, 4403 sent and read, a minute's back-off, one `error` log, the refusal in `channel.State` (criteria 1, 2, 4) |
| 1 | T-1248 | | Dashboard: protocol 2, flai's stamp checked before the proof, its own stamp in the answer, 4403, the refusal and the signed mark kept per project (criteria 1, 4) |
| 2 | T-1249 | T-1247 | `flai serve status` and `flai host status` print the refusal (criterion 2) |
| 2 | T-1250 | T-1248 | The banner and the header badge show a refused flai, as they show one that lacks methods (criterion 3) |
| 2 | T-1251 | T-1248 | `/api/projects`, the switcher, and Settings show signed, unsigned, or refused per project (criterion 3) |
| 3 | T-1252 | T-1249, T-1250, T-1251 | Operator, user, and design docs (criterion 5) |

## Figures

- Forecast: 1h30m, kept over flai's 1h1m. S-0235 has half this story's code and was planned at 1h.
- Cost of delay: 3.95 USD a week, from `flai cod`, kept.
- Touches: 12 files added from the layout and the design. The two declared folders are kept. The story's Notes record each touch.

## Assumptions

- S-0235 adds a check of a build's own stamp in `flai/internal/buildinfo/stamp.go` and `flaiover/src/lib/server/release.ts`. This story adds a check of a peer's stamp beside it, in the same two files. Those files are S-0235's, and S-0237 already waits for S-0235.
- A protocol 1 peer gets 4403, not a protocol error. The dashboard reads an old `hello` far enough to judge its stamp.
- `flai host status` reads `flai serve`'s state from the command. `host.Status` stays as it is unless the agent finds that cannot work.

## Risk: publishing before S-0239

Once S-0237 is released, a release flai refuses an unsigned dashboard, and a release dashboard refuses an unsigned flai. This repository runs both unsigned: `scripts/flai.sh` and `flai dashboard --build`. `dashboard.allow_unsigned` arrives only with S-0239, which waits for S-0237 and S-0238.

So publishing S-0237 before S-0239 would cut this repository's own dashboard off from its host flai. That lasts until S-0239 is released.

I recommend publishing S-0237 and S-0239 in one release. Publishing is already a deliberate, batched step (ADR-0032), so no task changes for this. The other way is to move the allowance into S-0237, but that merges part of S-0239 and I do not recommend it.

## Proposals

None. I would not split, merge, or drop any task.

### 2026-10-07T23:02:15Z orchestrator
Approved. The six tasks cover the five criteria:

- T-1247 and T-1248 cover criterion 1.
- T-1247 and T-1249 cover criterion 2.
- T-1250 and T-1251 cover criterion 3.
- T-1247 and T-1248 cover criterion 4, with their tests.
- T-1252 covers criterion 5.

No two tasks of a layer share a path. Both folder touches narrow to the files the tasks name. The 1h30m forecast and the 3.95 USD a week value stand.

On the release risk: I take your recommendation. The release policy is judgement, and publishing is mine while `publish` is on. Under it, I will not publish S-0237 before S-0239 is accepted, so `dashboard.allow_unsigned` reaches the remote in the same release (ADR-0032 keeps publishing a batched step). The operator can make the hold mechanical by setting `orchestration.release.whole_epics`. That would hold all of E-0015 until its last story.

Source: design/adrs/0032-accepting-a-story-merges-it-publishing-is-a-deliberate-batched-step-over.md

### 2026-10-07T23:02:16Z orchestrator
Resolved: Plan approved by the orchestrator under plan_backlog_stories; S-0237 will not be published before S-0239
