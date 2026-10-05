---
id: TH-0134
title: "S-0221 plan: seven tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0221-the-orchestrator-accepts-stories-in-review-when-permitted.md
  item: S-0221
status: resolved
participants: [planner-S-0221, alex]
created: 2026-10-05T04:48:53Z
updated: 2026-10-05T05:03:19Z
---

# TH-0134 S-0221 plan: seven tasks in four layers

On wip/kanban/stories/S-0221-the-orchestrator-accepts-stories-in-review-when-permitted.md.

## Entries

### 2026-10-05T04:48:53Z planner-S-0221
Plan for S-0221. All seven tasks are in the backlog.

Layers:

- Layer 1
  - T-0898: the ADR refining ADR-0032, plus `workflow.md` and `strategic-agents.md`. It fixes the conditions and the evidence format.
- Layer 2, all after T-0898; they run together and share no path
  - T-0900: `preview.Accept` gains an orchestrator option. It adds four blockers: the verified commit is not the head, a criterion is unticked, a file is outside the touches, or a thread is open.
  - T-0903: the `flai guard` rule for `accept_reviews`, and `item_move`'s refusal naming `flai accept`.
  - T-0905: the orchestrator's prompt, its agent file, and the convention's "As the orchestrator", template first.
- Layer 3
  - T-0907, after T-0900 and T-0903: `flai accept --by orchestrator` and `accept.run`. They check the permission, refuse on any blocker, record `by: orchestrator`, and write the evidence to the story's Notes. Its tests cover an acceptance, a refusal on an open thread, and the permission off.
- Layer 4, both after T-0907
  - T-0909: the review page and the story page show who accepted.
  - T-0911: `flai-cli.md`, `flaiover-dashboard.md`, and the users' guides.

Figures: touches grew from 13 to 25, keeping every declared one. The forecast is 1h30m, delivered 2026-10-05T14:23Z. The cost of delay is 88.82 USD/week, the story's share of E-0016's figure. The story's Notes, under Planning, say why.

Assumptions:

1. How flai knows "the verifier's run passed". Nothing records a verifier's verdict today: close-out prints its result and stores nothing. My recommendation:
   - The orchestrator runs its own verifier sub-agent on the story's worktree.
   - It passes the commit the verifier checked, and the evidence, to `flai accept`.
   - flai refuses when that commit is not the branch head.

   The alternative is for `scripts/close-out.sh` to record a pass on the story, which flai then checks. That widens the story into the close-out script and the template. Say if you prefer it.
2. S-0218 supplies `orchestration.permissions.accept_reviews`, the orchestrator's identity `orchestrator`, its prompt in `flai/internal/harness`, its agent file, and the guard's permission table with logged refusals. The tasks extend these and do not create them.
3. `item_move` keeps refusing a move to done, the orchestrator included. The orchestrator accepts through `flai accept`, so every acceptance it makes runs the preview and records the evidence.
4. "The review page shows who accepted" means the message after an acceptance, plus a line on the done story's page. The history list already shows each transition's `by`.
5. S-0276 changes `preview.Accept` and is ahead in the pull order. T-0900 starts from what S-0276 leaves on main.

Nothing to split, merge, or drop.

### 2026-10-05T05:03:19Z alex
Resolved.
