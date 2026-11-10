---
id: TH-1036
title: "Plan for S-0694: build ADR-0254 §5a's send-back and kluster blocker, six tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0694-adr-0254-5a-s-send-back-of-a-story-whose-kluster-deploy-or-checks-failed-and.md
  item: S-0694
status: resolved
participants: [planner-S-0694, orchestrator]
created: 2026-11-10T10:34:15Z
updated: 2026-11-10T11:10:08Z
---

# TH-1036 Plan for S-0694: build ADR-0254 §5a's send-back and kluster blocker, six tasks in four layers

On wip/kanban/stories/S-0694-adr-0254-5a-s-send-back-of-a-story-whose-kluster-deploy-or-checks-failed-and.md.

## Entries

### 2026-11-10T10:34:15Z planner-S-0694
**Recommended:** approve the plan as written. The solution is ADR-0254 §5a, which you settled on TH-0832 (D-0481). The plan needs no new decision from you.

## Tasks and layers

| Layer | Task | What it builds |
|-------|------|----------------|
| 1 | T-3024 | `kluster.Record` keeps the checks outcome and commit; an exported gate answers whether acceptance may go on |
| 2 | T-3025 | flai serve sends a failed story back to in-progress as `flai`, starts its agent, and keeps its namespace until it is cancelled or accepted |
| 2 | T-3028 | the `kluster` blocker in `flai accept` (operator, orchestrator, `--batch`) and its previews; `--without-kluster`, operator only, recorded in Notes |
| 2 | T-3029 | the merge queue lists a refused story with the reason `kluster`, so no batch is built and then dropped at landing |
| 3 | T-3035 | kluster.md, workflow.md, strategic-agents.md, flai-cli.md, and the user guide; the reference and the settings flag table are regenerated |
| 4 | T-3038 | close I-0268 |

Layer 2's three tasks share no file, so they can run together. Each one calls T-3024's gate.

## Figures

- Forecast: 2h, delivered 2026-11-10T16:10Z. flai said 2m from 2 criteria and 0 touches. S-0579, the pass this story extends, took 3033 s.
- Cost of delay: 50 USD a week, from your 20m per cycle. Kept as flai worked it out.
- Touches: 26 files, no folder touches, no `after`.

## Assumptions

1. The head that counts is the story branch's head, `refs/heads/story/<ID>`. That is the commit `KlusterEnv` deploys. A batch's rebased head is not deployed, so it is not the head that counts.
2. The checks outcome is recorded in the kluster record (`flai/internal/kluster`), not in serve's per-story `ChecksRun`, so `flai accept` can read it outside package serve.
3. A sent-back story keeps its namespace while it is in progress. Cancelling or accepting it tears the namespace down. That follows §5a's second bullet.
4. `--without-kluster` is refused with `--batch` and `--by orchestrator`, as `--without-queue-tiers` is, since §5a makes it "the operator's word only".
5. The queue's eligibility reason is an addition: §5a says only that the queue's acceptance is refused. Without the reason, the queue would verify a batch and then drop it at landing.

No split, merge, or drop is proposed.

### 2026-11-10T11:09:25Z orchestrator
Approved.

- The solution is ADR-0254 §5a as the operator settled it on TH-0832 (D-0481), so it needs no new decision.
- T-3024's gate comes first. The send-back (T-3025), the blocker (T-3028), and the queue's reason (T-3029) share no file and run together. Then come the docs (T-3035) and the close (T-3038). Together they cover both criteria.
- The queue's `kluster` reason (assumption 5) keeps §5a's refusal from being found only at landing, so it is in scope.
- `--without-kluster` is refused to `--batch` and `--by orchestrator`, which keeps it the operator's word, as §5a says.
- There are 26 file touches and no folder. The 2h forecast states its reason. The value is 50 USD a week, from the operator's 20m input.

I am finalizing S-0694.

Source: design/conventions/strategic-agents.md § As the orchestrator

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-11-10T11:10:08Z orchestrator
Resolved: Plan approved; S-0694 finalized
