---
id: TH-0595
title: "S-0472 duplicates S-0477: cancel it once S-0477 is accepted?"
anchor:
  path: wip/kanban/stories/S-0472-flai-guard-refuses-a-story-s-sub-agents-flai-components-list-and-serve-which-only-read-or-run-the-project-s-own-dev-server.md
  item: S-0472
status: resolved
participants: [agent-S-0477, orchestrator, agent-S-0472, alex]
created: 2026-10-22T21:00:46Z
updated: 2026-10-23T17:53:26Z
---

# TH-0595 S-0472 duplicates S-0477: cancel it once S-0477 is accepted?

On wip/kanban/stories/S-0472-flai-guard-refuses-a-story-s-sub-agents-flai-components-list-and-serve-which-only-read-or-run-the-project-s-own-dev-server.md.

## Entries

### 2026-10-22T21:00:46Z agent-S-0477
I-0156 (S-0472) and I-0165 (S-0477) record the same refusal: flai guard has no `components` entry, so a story's sub-agent cannot run `flai components list` or `serve`. S-0477 fixes it: list becomes a read for every agent, serve is allowed to a story's task sub-agent (not a read-only one, not the strategic agents). Recommendation: cancel S-0472 and close I-0156 as fixed by S-0477, so that no agent is started on it once S-0477 moves to review. S-0477 does not wait on this answer.

### 2026-10-22T21:00:54Z orchestrator (recommendation)
Recommended for the operator: cancel S-0472 now, with the reason "duplicate of S-0477, which fixes I-0156's refusal". Then answer TH-0593's merge so that S-0477 closes I-0156 too.

- The two stories fix the same missing `components` entry in the guard. S-0477 is in progress with the fix.
- S-0472 is in ready. Today it is held on overlap while S-0477 is in progress. Once S-0477 moves to review, that hold lifts, and an agent could start S-0472 on a change already made. Cancelling before then avoids that.
- If you would rather keep S-0472 for closing I-0156 alone, give it `after: S-0477`, so that it waits for the fix to land.

This is the same question as TH-0593. Answer either one and I will treat both as settled.

Source: design/conventions/work-management.md § Rules

### 2026-10-22T21:32:33Z orchestrator (recommendation)
Update for the operator: S-0477 is in review, verified at e54d72a7. It closes I-0165 and, by ADR-0186, lets every agent run `flai components list` and a story's task sub-agent run `serve`. That fixes I-0156's refusal too. I am resolving TH-0593 so that the merge queue can take S-0477. This thread now carries the S-0472 question alone.

Recommended: cancel S-0472 as a duplicate of S-0477, and close I-0156 with `flai issue close I-0156 --reason "fixed by S-0477 (ADR-0186)"`. S-0477 no longer closes I-0156 in its own diff, but nothing else is lost. If you would rather keep S-0472, give it `after: S-0477`, and its agent closes I-0156 once S-0477 lands.

Source: design/conventions/work-management.md § Rules

### 2026-10-23T06:22:50Z agent-S-0472
S-0472's agent here. flai serve started S-0472 after S-0477 was accepted. So I took the thread's second option: S-0472 adds no code and only closes I-0156, as fixed by S-0477 (ADR-0186), with T-2067. S-0477's guard test `TestASubAgentListsComponentsAndATaskSubAgentServesThem` reproduces I-0156's refusal and passes on S-0472's branch. S-0472 goes to review with that one commit. Accepting it lands I-0156's close. Cancelling it leaves I-0156 open on main, to close by hand.

### 2026-10-23T06:24:00Z orchestrator (recommendation)
Recommended for the operator: resolve this thread. My earlier recommendation, to cancel S-0472, is out of date.

- **S-0472 took the second option.** It adds no code and only closes I-0156 as fixed by S-0477 (ADR-0186).
- **It is verified at 1902430c.** The verifier confirmed that the fix and its test are on main, and that both criteria are met.
- **This thread is now the only thing holding it.** The merge queue refused S-0472 with `thread_open: thread TH-0595`. agent-S-0477 opened the thread, so I cannot resolve it.

Once it is resolved, I enqueue S-0472 and it lands with I-0156's close.

Source: design/conventions/strategic-agents.md § As the orchestrator

### 2026-10-23T17:53:26Z alex
Resolved: S-0472 was accepted
