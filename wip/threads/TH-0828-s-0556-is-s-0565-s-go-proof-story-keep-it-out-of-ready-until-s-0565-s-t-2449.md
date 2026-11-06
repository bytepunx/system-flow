---
id: TH-0828
title: "S-0556 is S-0565's Go proof story: keep it out of ready until S-0565's T-2449 sets its harness to container, as TH-0826 asks for S-0558"
anchor:
  path: wip/kanban/stories/S-0556-flai-guard-still-lists-the-verify-mcp-tool-among-a-sub-agent-s-reads-though-it.md
  item: S-0556
status: open
participants: [orchestrator]
created: 2026-11-05T01:44:20Z
updated: 2026-11-05T01:44:20Z
---

# TH-0828 S-0556 is S-0565's Go proof story: keep it out of ready until S-0565's T-2449 sets its harness to container, as TH-0826 asks for S-0558

On wip/kanban/stories/S-0556-flai-guard-still-lists-the-verify-mcp-tool-among-a-sub-agent-s-reads-though-it.md.

## Entries

### 2026-11-05T01:44:20Z orchestrator
**Recommended:** decide S-0556 the same way you decide S-0558 on TH-0826. My recommendation there is to make the story a draft again, so that it cannot be promoted before S-0565 uses it:

```bash
flai edit S-0556 --draft
```

**What happened:**

- planner-S-0556 planned S-0556. I approved the plan on TH-0827 and set `time_lost_per_cycle: 5m`. The value is 12.50 USD a week and the forecast is 40m.
- `flai promote --drafts` then listed S-0556 as complete. TH-0821's confirmed answer says to plan and finalize both proof stories before S-0565 starts. So I finalized it under `finalize_drafts`.
- `flai promote --candidates` now lists it at rank 2 (40m, throughput), behind S-0558.

**Why I have not promoted it:**

- TH-0821 makes S-0556 S-0565's Go proof story for criterion 3. S-0565's T-2449 sets its harness to `container`, and only after that does it move to ready (S-0565's assumption 4).
- If it goes to ready now, an ordinary claude-code agent works it outside the container, and the proof is used up.
- S-0565 is still a draft, behind S-0520 and the host preconditions.

**Why I am asking you:** flai guard refuses the orchestrator `draft: true`. Until you decide, I skip S-0556 in every promotion while this thread is open.

**Alternative:** leave it finalized. S-0565's agent or you then move it to ready once its harness is set. That relies on every promotion run seeing this thread.
