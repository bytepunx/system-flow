---
id: TH-0201
title: "S-0300: criteria 1 and 2 are already met; what does \"draft tasks\" in criterion 3 mean?"
anchor:
  path: wip/kanban/stories/S-0300-the-planner-agent-should-be-available-on-epic-work-items.md
  item: S-0300
status: open
participants: [planner-S-0300]
created: 2026-10-06T21:46:30Z
updated: 2026-10-06T21:48:26Z
---

# TH-0201 S-0300: criteria 1 and 2 are already met; what does "draft tasks" in criterion 3 mean?

On wip/kanban/stories/S-0300-the-planner-agent-should-be-available-on-epic-work-items.md.

## Entries

### 2026-10-06T21:46:30Z planner-S-0300
Two of S-0300's criteria are already met in the code. One word in the third needs your call.

**Criteria 1 and 2 are met already.**

- **Epic page:** `PlanAction.svelte` shows **Plan** on an open epic's page while the `plan` host action is on (S-0208, `flaiover/src/routes/items/[id]/+page.svelte:432`, where `canEdit` includes epics).
- **Card menu:** `cardMenu` offers **Plan** on an open epic card (S-0263, `flaiover/src/lib/cardmenu.ts:34`, tested at `cardmenu.test.ts:71`). It has been released since flai 1.31.1.
- **Board view:** epic cards appear on the board only when the **epics** checkbox is ticked. That may be why you didn't see Plan there.

I plan one small task to cover these: a test that the item page shows Plan on an open epic, since no item-page test has one now.

**Criterion 3 needs your call.** Planning an epic drafts its stories today, but not their tasks. Tasks have no draft flag.

**Recommended:** when the planner drafts a story for an epic, it also drafts that story's tasks in the same run. It enriches each drafted story as it would a story (touches, forecast, cost of delay value), then creates the tasks in the backlog with `## Work`, `## Done when`, touches, and `after`. The tasks get no draft flag of their own. They count as drafts because their story is one, and the story's agent reviews them when it pulls the story. This changes the planner's epic prompt (`flai/internal/harness/harness.go`), the convention, the design, the docs, and both `planner.md` agent definitions. No schema change and no ADR.

**Alternatives:**

- **(b)** Add a `draft` flag on tasks, with flai check, the dashboard, and finalizing to match. That needs an ADR and is several stories' work.
- **(c)** Keep one epic run for the stories. flai serve then queues a story planner run for each story the epic run drafted. That costs more runs and more code in `serve`, but each run's context stays small.

I am drafting S-0300's tasks on the recommendation now and will adjust them to your answer.

### 2026-10-06T21:48:26Z planner-S-0300
The plan for S-0300 is written, based on the recommended answer above.

**Tasks and layers.**

- **Layer 1** (these run together and share no path):
  - T-1047: the design (`design/system/strategic-agents.md`) and the convention, plus its template copy and the CHANGELOG, say that an epic's planner drafts the tasks of each story it drafts.
  - T-1048: `planPrompt`'s epic text in `flai/internal/harness/harness.go`, with its tests. The planner enriches each story it drafts and drafts that story's tasks, and its thread and summary name them.
  - T-1049: an item page test that shows **Plan** on an open epic and not on a done or archived one. This covers criteria 1 and 2, which are already met.
- **Layer 2** (each waits for T-1047, so its words follow the convention's):
  - T-1050: both `planner.md` agent definitions. A `.claude/` write needs your **allow**.
  - T-1051: `docs/users/flaiover.md`, `docs/users/flai.md`, and `docs/operators/index.md`.

**Figures.**

- **Forecast:** 40m, up from flai's 24m for the copies, the permission wait, and both test suites.
- **Delivery:** 2026-10-07T00:30Z. S-0300 is held behind S-0229 and S-0299 on `harness.go`.
- **Cost of delay:** 600 USD a week, from your 4h a cycle.

**Assumptions.**

- Tasks need no draft flag.
- No ADR is needed: the change extends what ADR-0082 and S-0209 give the epic planner.
- `flai plan` and the host API stay unchanged.

**Proposals.**

1. Drop the declared touch `flai/cmd`: no task changes it, and while the story is in progress it holds every ready story under `flai/cmd`. Narrowing `flaiover/src` is not needed, because T-1049's file replaces it in the claim. I have kept both touches, since they are yours.
2. If you choose (c) instead, T-1048 moves to `flai/internal/serve`, and T-1050 and T-1051 change their wording. If you choose (b), S-0300 should be split, with the task draft flag as a story of its own.
