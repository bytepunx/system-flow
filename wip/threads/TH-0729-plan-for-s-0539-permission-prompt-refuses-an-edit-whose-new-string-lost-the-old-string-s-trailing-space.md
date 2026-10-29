---
id: TH-0729
title: "Plan for S-0539: permission_prompt refuses an edit whose new string lost the old string's trailing space"
anchor:
  path: wip/kanban/stories/S-0539-an-edit-whose-old-and-new-strings-end-in-a-space-drops-the-space-so-a-protected-claude-settings-json-write-needs-a-second-permission-round-to-fix.md
  item: S-0539
status: open
participants: [planner-S-0539]
created: 2026-10-29T03:03:44Z
updated: 2026-10-29T03:03:44Z
---

# TH-0729 Plan for S-0539: permission_prompt refuses an edit whose new string lost the old string's trailing space

On wip/kanban/stories/S-0539-an-edit-whose-old-and-new-strings-end-in-a-space-drops-the-space-so-a-protected-claude-settings-json-write-needs-a-second-permission-round-to-fix.md.

## Entries

### 2026-10-29T03:03:44Z planner-S-0539
S-0539 is planned with three tasks in two layers. No input of yours is missing.

## Cause

TH-0712 shows what `permission_prompt` received from S-0416's agent: an old string ending in `exit 0;` and a space, and a new string ending in `exit 0;` alone. Claude Code drops trailing whitespace from an Edit's new string before flai sees the call. flai hands the input back unchanged, so it cannot restore the space.

## Proposed remedy

`permission_prompt` refuses at once, with no thread and under auto-approve too, an `Edit` or a `MultiEdit` edit of a protected path whose old string ends in a space or a tab while its new string does not. The refusal tells the agent to end both strings on a non-space character. The agent retries at once instead of writing the wrong text and asking you a second time.

## Tasks

| Layer | Task | after | Touches |
|-------|------|-------|---------|
| 1 | T-2320 the refusal, with a test reproducing S-0416's instance | none | `flai/internal/mcpserver/permission.go`, `flai/internal/mcpserver/permission_test.go` |
| 2 | T-2321 document it | T-2320 | `docs/users/flai.md`, `design/system/agent-adapters.md` |
| 2 | T-2322 close I-0223 | T-2320 | I-0223's file, `design/issues/summary.md` |

T-2321 and T-2322 share no path and run together.

## Figures

- Forecast: 1h, up from flai's 3m, which counted no touches. Delivery stays at 2026-10-29T22:24:00Z.
- Cost of delay value: 37.50 USD a week, from flai's input of 15m lost per cycle.

## Assumptions

- The cause is Claude Code's trimming, not the model writing the new string without the space. T-2320 confirms this from TH-0712, TH-0713, and TH-0720 first, and asks here if they show otherwise.
- The check covers protected writes only, which is where I-0223's cost is. An ordinary file loses the space the same way, but a second plain edit fixes it at no cost to you. I would not extend the check to `flai guard` for every Edit unless the issue recurs there.
- No convention line is added. The refusal tells the agent what to do at the moment it matters, so no agent's prime carries the rule.
- S-0532, in progress, claims `flai/internal/mcpserver` and `docs`. S-0539 is held until S-0532 moves to review. It names no story in `after`.
