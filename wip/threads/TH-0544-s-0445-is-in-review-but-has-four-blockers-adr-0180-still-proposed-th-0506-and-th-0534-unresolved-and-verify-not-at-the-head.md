---
id: TH-0544
title: "S-0445 is in review but has four blockers: ADR-0180 still proposed, TH-0506 and TH-0534 unresolved, and verify not at the head"
anchor:
  path: wip/kanban/stories/S-0445-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md
  item: S-0445
status: open
participants: [orchestrator]
created: 2026-10-21T00:52:09Z
updated: 2026-10-21T00:52:09Z
---

# TH-0544 S-0445 is in review but has four blockers: ADR-0180 still proposed, TH-0506 and TH-0534 unresolved, and verify not at the head

On wip/kanban/stories/S-0445-the-host-s-installed-flai-is-again-behind-the-project-s-releases-after-i-0149-closed-so-the-mcp-task-done-refuses-every-close.md.

## Entries

### 2026-10-21T00:52:09Z orchestrator
**Recommended:** S-0445's agent does the following, and then I verify again and accept S-0445:
1. Run `flai adr accept 180` on the branch. The operator confirmed option B on TH-0506, which is the decision ADR-0180 records.
2. Resolve TH-0506 and TH-0534. Option B is built, and the close-out now passes smoke.
3. Commit, and let the close-out's verify pass at the new head.

## The blockers

- **ADR-0180 is `proposed`.** `design/adrs/0180-flai-task-done-refuses-an-older-flai-only-when-the-newest-flai-release-raised.md` records the task-done rules number, the decision both TH-0506 and the approved plan (TH-0531) rest on. A proposed ADR would land the rule with no accepted decision behind it.
- **TH-0506 is answered, not resolved.** The dry run says so. It is the agent's own thread.
- **TH-0534 is open, not resolved.** The operator confirmed my recommendation there, but the agent has not resolved the thread.
- **The verify is not at the head.** The last `flai verify S-0445` passed at 48e1074df80d, but story/S-0445 is now at df83adf5053b. The run will be repeated after the ADR commit in any case.

## What already passes

The last verify passed every step, including integration and smoke. The branch contains main, and every task is done or cancelled. I'll run the verifier on the criteria once the head is final.
