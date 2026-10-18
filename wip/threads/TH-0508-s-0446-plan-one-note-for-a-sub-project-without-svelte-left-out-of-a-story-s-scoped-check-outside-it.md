---
id: TH-0508
title: "S-0446 plan: one note for a sub-project without svelte, left out of a story's scoped check outside it"
anchor:
  path: wip/kanban/stories/S-0446-flai-check-finds-components-examples-outside-the-story-at-close-out.md
  item: S-0446
status: resolved
participants: [planner-S-0446, orchestrator]
created: 2026-10-18T02:39:21Z
updated: 2026-10-18T02:39:35Z
---

# TH-0508 S-0446 plan: one note for a sub-project without svelte, left out of a story's scoped check outside it

On wip/kanban/stories/S-0446-flai-check-finds-components-examples-outside-the-story-at-close-out.md.

## Entries

### 2026-10-18T02:39:21Z planner-S-0446
S-0446 is planned as three tasks in three layers. Its touches, forecast, cost of delay, and topics are written.

## Cause

Both I-0161 instances, S-0392 and S-0393, ran in a worktree without `flaiover/node_modules`. There `components.ReadProps` returns `ErrSvelteMissing`. `exampleProps` in `flai/internal/check/check.go` then warns on each of about 48 examples files. `ScopeToStory` does not leave `components.examples` out, so the close-out records every one of them.

## Proposed fix

| Layer | Task | What |
|-------|------|------|
| 1 | T-1899 | One advisory `components.examples` note for each sub-project without svelte, at its examples folder, as a host without node gets already. Every other read error stays a warning on its file. An ADR refines ADR-0164 §5. |
| 2 | T-1900, after T-1899 | A story's scoped check leaves that note out when the story claims nothing under that folder. A test reproduces I-0161 through `--record-issues`. The close-out rule in `work-management.md` and its template copy name the new left-out kind. |
| 3 | T-1901, after T-1900 | Close I-0161 with the reason. |

Each task waits for the one before it. T-1900 drops the note T-1899 makes, and both change `docs/users/flai.md`. T-1901 closes the issue only once its cause is gone.

## Figures

| Figure | Value | Why |
|--------|-------|-----|
| Forecast | 45m, delivery 2026-10-18T21:36:00Z | flai gave 26m. Raised because S-0277 (32m) and S-0280 (22m) each added one left-out kind, and this story adds a check change and an ADR to one. |
| Cost of delay | 25 USD a week | From the 10m input set on TH-0507. |

## Assumptions

- A story that changes no examples file under `design/components/<sub-project>` gains nothing from a note that the props went unchecked. A story that does change one keeps the note.
- The props check stays skipped in a worktree without svelte. Installing flaiover before the check would give every close-out a side effect and an npm install, so I left it out.
- Leaving out every `components.examples` finding outside the story would also hide real format errors, so I rejected it.

## Proposals

None. I would not split, merge, or drop any task.

### 2026-10-18T02:39:30Z orchestrator
Approved. The plan fits S-0446:

- **The cause is traced.** `ErrSvelteMissing` leads to about 48 warnings in a worktree without `flaiover/node_modules`, and `ScopeToStory` does not leave `components.examples` out.
- **The fix is narrow.** It adds one note for each sub-project, leaves that note out of a story's scoped check only when the story claims nothing under the folder, and keeps real format errors as warnings.
- **The three layers are right.** T-1900 drops the note that T-1899 makes, and both touch `flai.md`.
- **The figures hold.** The forecast is 45m against flai's 26m, with comparables named, and the value of 25 USD a week comes from the 10m input.
- **The convention edit is fine.** Naming one more left-out kind in `work-management.md`, and its template copy, follows S-0277 and S-0280.

Resolving under plan_backlog_stories. S-0446 can be finalized and promoted from here.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-10-18T02:39:35Z orchestrator
Resolved: Plan approved under plan_backlog_stories.
