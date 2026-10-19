---
id: TH-0517
title: S-0449 is verified and its criteria are met, but TH-0515 is answered and not resolved, which blocks acceptance
anchor:
  path: wip/kanban/stories/S-0449-improvements-to-design-components-screen-part-1.md
  item: S-0449
status: resolved
participants: [orchestrator, alex]
created: 2026-10-19T14:36:24Z
updated: 2026-10-19T17:30:21Z
---

# TH-0517 S-0449 is verified and its criteria are met, but TH-0515 is answered and not resolved, which blocks acceptance

On wip/kanban/stories/S-0449-improvements-to-design-components-screen-part-1.md.

## Entries

### 2026-10-19T14:36:24Z orchestrator
**Recommended:** S-0449's agent, or you, resolves TH-0515 with "settled by the operator's answer: (A) with a one-element array, built in T-1907 and T-1908". I'll then accept S-0449 at 572deff2e2a0, or you can accept it yourself.

## The one blocker

`flai accept S-0449 --by orchestrator --verified 572deff2e2a0 --dry-run` says: "blocked: thread TH-0515 on S-0449 is answered, not resolved". agent-S-0449 opened that thread, and I may not resolve a thread I did not open.

## Everything else is ready

- **Verify:** `flai verify S-0449` passed every step at the head, 572deff2e2a0, against main 882a129, including integration and smoke. The branch contains main.
- **Criteria:** the verifier matched all five to changed files, and every box is ticked.
  - 1: the harness's 10px padding.
  - 2: whole-component measurement and the frame's scroll bars.
  - 3: skeletons resolved with the sub-project's `typescript`, filled and sent.
  - 4: the select-only mount, with stale messages dropped.
  - 5: ADR-0172's 501 answers for the app's own requests, and a headless pass with no 500.
- **Design and docs:** component-library.md, flaiover-dashboard.md and flai-cli.md record the changed rules. The user guides and the reference are updated.

## A note on criterion 5's "file these as issues"

The 500s S-0444's pass found were fixed in this story, not filed, and the pass afterwards saw none. The story filed I-0164 to I-0166 for what it found besides. I count fixing the 500s as meeting the criterion. If you want the original 500 recorded as an issue too, say so here.

### 2026-10-19T17:30:21Z alex
Resolved: S-0449 was accepted
