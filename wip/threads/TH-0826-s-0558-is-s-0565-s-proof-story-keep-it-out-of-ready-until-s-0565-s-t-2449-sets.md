---
id: TH-0826
title: "S-0558 is S-0565's proof story: keep it out of ready until S-0565's T-2449 sets its harness to container; recommend making it a draft again"
anchor:
  path: wip/kanban/stories/S-0558-safety-md-says-every-force-push-needs-confirmation-while-git-md-exempts-flai.md
  item: S-0558
status: open
participants: [orchestrator]
created: 2026-11-05T01:37:16Z
updated: 2026-11-05T01:37:16Z
---

# TH-0826 S-0558 is S-0565's proof story: keep it out of ready until S-0565's T-2449 sets its harness to container; recommend making it a draft again

On wip/kanban/stories/S-0558-safety-md-says-every-force-push-needs-confirmation-while-git-md-exempts-flai.md.

## Entries

### 2026-11-05T01:37:16Z orchestrator
**Recommended:** make S-0558 a draft again, so it cannot be promoted before S-0565 is ready to use it as its proof:

```bash
flai edit S-0558 --draft
```

**What happened:**

- planner-S-0558 planned S-0558. I settled its two threads under `plan_backlog_stories`. TH-0824 set `time_lost_per_cycle: 3m`, and the value is 7.50 USD a week. TH-0825 approved the plan.
- `flai promote --drafts` then listed S-0558 as complete, so I finalized it under `finalize_drafts`.
- `flai promote --candidates` now lists it at rank 1 (12m, throughput).

**Why I have not promoted it:**

- TH-0821 settled S-0558 as S-0565's proof story for criteria 1 and 2. That answer was confirmed by flai serve after 48h. S-0565's T-2449 sets the proof story to `harness: container`, and only then is it moved to ready (S-0565's assumption 4).
- If S-0558 goes to ready now, an ordinary claude-code agent works it outside the container. That uses up the proof.
- S-0565 is still a draft. It waits behind S-0520 and the host preconditions.

**Why I am asking you:** my undo, `draft: true`, was refused by flai guard: the orchestrator may only finalize a draft, never make one. Until you decide, I will skip S-0558 in every promotion while this thread is open.

**Alternative:** leave S-0558 finalized, and have S-0565's agent or you move it to ready once its harness is set. That relies on every promotion run seeing this thread.
