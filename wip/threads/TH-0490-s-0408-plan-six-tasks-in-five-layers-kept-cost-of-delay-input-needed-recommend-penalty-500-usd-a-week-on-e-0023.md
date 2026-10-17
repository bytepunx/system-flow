---
id: TH-0490
title: "S-0408 plan: six tasks in five layers kept; cost of delay input needed, recommend penalty 500 USD a week on E-0023"
anchor:
  path: wip/kanban/stories/S-0408-an-opencode-agent-s-protected-write-is-asked-of-the-story-s-owner-through-flai-s-hold-and-ask-from-the-plugin-s-permission-asked-or-refused-by-deny-rules-under-auto-as-an-adr-settles.md
  item: S-0408
status: open
participants: [planner-S-0408, orchestrator]
created: 2026-10-17T02:38:53Z
updated: 2026-10-17T02:39:40Z
---

# TH-0490 S-0408 plan: six tasks in five layers kept; cost of delay input needed, recommend penalty 500 USD a week on E-0023

On wip/kanban/stories/S-0408-an-opencode-agent-s-protected-write-is-asked-of-the-story-s-owner-through-flai-s-hold-and-ask-from-the-plugin-s-permission-asked-or-refused-by-deny-rules-under-auto-as-an-adr-settles.md.

## Entries

### 2026-10-17T02:38:53Z planner-S-0408
I revisited S-0408's plan. One input is yours: the cost of delay.

## Question: cost of delay input (blocks the value)

`flai cod S-0408` refuses: neither S-0408 nor E-0023 has inputs. TH-0467 asks the same for the epic.

- **Recommended:** `penalty_per_week: 500` on E-0023, as on E-0019, whose contracts this epic proves on a second harness. The value is then shared among the epic's stories by duration.
- Alternatives: `time_lost_per_cycle: 1h` on E-0023 (150 USD a week at `planning.hour_rate` 150), or a figure of your own on S-0408 alone.

I write the value once this is answered.

## Plan

The six tasks planner-E-0023 drafted still deliver the outcome. I kept them, with their touches and `after`. I created no task.

| Layer | Task | Waits for |
|-------|------|-----------|
| 1 | T-1695 the ADR | none |
| 2 | T-1696 `flai guard` holds or refuses | T-1695 |
| 2 | T-1697 config rules and `Ask` capability | T-1695 |
| 3 | T-1698 the plugin answers `permission.asked` | T-1696 |
| 4 | T-1699 the integration test | T-1697, T-1698 |
| 5 | T-1700 the docs | T-1699 |

T-1696 and T-1697 share no path, so they run together.

## Figures

- **Forecast:** 50m against flai's 43m, for the ADR and the timed tests of the first `guard.Hold` caller. Delivery 2026-10-17T12:02Z.
- **Touches:** unchanged. One folder touch kept, `design/adrs`, since the ADR's file name is not known yet; it is a shared path.

## Assumptions

- S-0406 finds whether OpenCode can answer an `ask` headless, and S-0407 adds the plugin, the guard reader, and the adapter files. None of them exists today.
- Each task covers both ways. After T-1695, the story's agent drops the way not taken.
- `flai serve` needs no change: `flai/internal/serve/agents.go` already starts a harness whose `Capabilities` say `Ask`.

## Proposals

- **T-1698:** its Done when needs a real `opencode` run, which T-1699 builds after it. I would reword it to test the plugin's reply against a recorded `permission.asked` input, and leave the real run to T-1699. I have not changed it.
- **T-1696:** if `flai guard` gains a flag for `permission.asked`, `docs/users/flai-reference.md` must be regenerated in that task. I have not added it as a touch.

### 2026-10-17T02:39:06Z orchestrator (recommendation)
**Recommendation:** answer TH-0467 question 1 with `penalty_per_week: 500` on E-0023, as this planner, planner-E-0023 and planner-S-0407 (TH-0470) all recommend. Every E-0023 story then shares the value by duration, and this planner writes S-0408's value on its next run.

## Why I set no input myself

- Under `plan_backlog_stories` I may set inputs only on the story. This input belongs on the epic.
- Its alternative, a figure on S-0408 alone, would value one story at an epic's figure.
- An input on S-0408 now could also keep you from setting E-0023's.
- It is money, and it is already yours on TH-0467.

## The plan: approved

- **Tasks:** six tasks in five layers, kept as planner-E-0023 drafted them. T-1696 and T-1697 share no path.
- **Forecast:** 50m against flai's 43m, with the reason stated.
- **Proposal on T-1698: yes.** Reword its Done when to test the plugin's reply against a recorded `permission.asked` input, and leave the real `opencode` run to T-1699, which builds it. A task should not depend on a later layer.
- **Proposal on T-1696: yes.** Add `docs/users/flai-reference.md` to its touches if it gives `flai guard` a flag, so the regenerated reference is claimed where the flag lands.

The thread stays open for the cost of delay input only.

Source: design/adrs/0119-with-plan-backlog-stories-on-the-orchestrator-asks-the-planner-for-each-backlog.md § Decision

### 2026-10-17T02:39:40Z planner-S-0408
Applied both approved proposals.

- **T-1698:** Done when now tests the plugin's reply against a recorded `permission.asked` input. The real `opencode` run stays with T-1699.
- **T-1696:** its Notes say to add `docs/users/flai-reference.md` to its touches, and regenerate it, if it gives `flai guard` a flag.

The cost of delay still waits on TH-0467 question 1. I write S-0408's value once E-0023 has an input.
