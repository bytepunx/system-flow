---
id: TH-0164
title: "Plan for S-0284: let the project's owner answer a permission thread, record it in an ADR, document it, close I-0081"
anchor:
  path: wip/kanban/stories/S-0284-flai-s-permission-prompt-waits-for-an-answer-from-the-story-s-owner-and-the-operator-s-thread-replies-carry-another-name-so-their-allow-is-never-seen-and-the-write-times-out.md
  item: S-0284
status: resolved
participants: [planner-S-0284, alex]
created: 2026-10-06T06:25:19Z
updated: 2026-10-06T06:47:24Z
---

# TH-0164 Plan for S-0284: let the project's owner answer a permission thread, record it in an ADR, document it, close I-0081

On wip/kanban/stories/S-0284-flai-s-permission-prompt-waits-for-an-answer-from-the-story-s-owner-and-the-operator-s-thread-replies-carry-another-name-so-their-allow-is-never-seen-and-the-write-times-out.md.

## Entries

### 2026-10-06T06:25:19Z planner-S-0284
The plan for S-0284. It has four tasks in three layers.

**Proposed fix. Please confirm or redirect, since the story asks for a proposal before building.** I-0081's one instance: S-0218's owner is `arobson`, but you replied `allow` on TH-0158 as `alex`. That is the `owner` in `system-flow.yaml`, and the name the dashboard always writes as. `awaitAnswer` in `flai/internal/mcpserver/permission.go` takes an answer only from the story's owner, so the write waited until Claude Code's 1800 s MCP idle timeout ended it. Recommended fix: an answer may come from the story's owner **or** the project's owner (`repo.Manifest.Owner`). Entries by the asking agent or by other agents are still not answers. The request names both people when they differ. Another option is "anyone but an agent", but flai has no reliable way to tell a person's name from an agent's, so I recommend against it.

Tasks:

1. Layer 1: **T-0997**, an ADR that refines ADR-0086's decision 2 on who answers. It starts by taking your answer here.
2. Layer 2, after T-0997: **T-0998**, the change to `awaitAnswer`, `askOperator`, and `permissionRequest`, with a test that reproduces I-0081 (story owner `arobson`, manifest owner `alex`, `alex` replies allow).
3. Layer 2, after T-0997: **T-0999**, `design/system/flai-cli.md` and `docs/users/flai.md` say who answers. It shares no paths with T-0998, so the two can run in parallel.
4. Layer 3, after T-0998 and T-0999: **T-1000**, close I-0081 with `flai issue close` and its reason.

Figures: forecast 30m, delivery 2026-10-06T10:45Z (flai's 18m plus 12m for the ADR and the proposal step). Cost of delay 75 USD a week, from your 30m-per-cycle input, unchanged.

Assumptions:

- Changing who may answer contradicts ADR-0086's accepted decision ("entries by anyone else … are not answers"), so it needs a new ADR, not an edit.
- `delegation.md` and its template copy say only "asks the operator" and stay as they are, so the template is not touched.
- S-0283 (I-0082) also changes `permission.go`. flai's overlap hold runs the two one after another, so I set no `after` between them. Whichever goes second merges the other's change.

I propose no splits, merges, or drops.

### 2026-10-06T06:47:24Z alex
Resolved.
