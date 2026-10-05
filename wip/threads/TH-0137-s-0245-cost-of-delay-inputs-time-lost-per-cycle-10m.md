---
id: TH-0137
title: "S-0245 cost of delay inputs: time lost per cycle 10m?"
anchor:
  path: wip/kanban/stories/S-0245-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md
  item: S-0245
status: resolved
participants: [planner-S-0245, alex]
created: 2026-10-05T05:43:43Z
updated: 2026-10-05T05:46:11Z
---

# TH-0137 S-0245 cost of delay inputs: time lost per cycle 10m?

On wip/kanban/stories/S-0245-flai-adr-new-numbers-from-the-story-s-worktree-only-so-parallel-story-branches-take-the-same-adr-number.md.

## Entries

### 2026-10-05T05:43:43Z planner-S-0245
S-0245 has no cost of delay inputs and no epic, so `flai cod S-0245` refuses and I cannot record a value. The inputs are yours to give.

**Recommended:** `time_lost_per_cycle: 10m`. That is the cost recorded on I-0063, the issue S-0245 remediates. It is also the figure you gave S-0252 (TH-0110), the twin fix for issue numbers. Each collision means renumbering an ADR at the rebase and fixing its links by hand: I-0063 has five instances in one day, across S-0200, S-0201, and S-0204. With `hour_rate` 150 and the default 168h cycle, 10m works out to about **25 USD/week**, the same as S-0252, so the pull order compares the two fairly.

**Alternative:** `time_lost_per_cycle: 20m`. That allows for the eight to ten links fixed by hand in S-0200's instances (TH-0084, TH-0085), which is about 50 USD/week.

Reply with the figure you want, or with "go with the recommendation". Meanwhile I am planning the touches, the forecast, and the tasks, which do not depend on the answer.

### 2026-10-05T05:44:25Z alex
20m

### 2026-10-05T05:45:24Z planner-S-0245
Thank you. I recorded `time_lost_per_cycle: 20m`, and `flai cod` works it out to **50.00 USD a week**, which is now the value. Here is the plan for S-0245.

**Proposed fix.** `adr.NextNumber` reads only the checkout's own `design/adrs`. With the fix, it also takes the names `storygit.FolderNames` returns for that folder. S-0252 added that helper for issue numbers, and it reads the main checkout, every worktree (uncommitted files included), the main branch, and every `story/*` branch. The next ADR is then one past the highest anywhere.

**Tasks, in two layers:**

- Layer 1: **T-0915**, flai adr new numbers past every ADR on main, every story worktree, and every story branch. This is the code, a reproduction test in `flai/cmd/adr_test.go` modelled on S-0252's, the help text, and the reference regeneration.
- Layer 2: **T-0916**, the design and the user guides say how ADRs are numbered, and I-0063 is closed. This covers `flai-cli.md`, `documentation-standard.md`, `flai.md`, and `flaiover.md`. It waits for T-0915 because it describes what that task built, and the issue closes only once the test passes.

**Assumptions:**

- I-0063's issue-number instance (S-0201, I-0062) is already fixed by S-0252, so this story covers ADRs only.
- The checks of `--supersedes` and `--refines` keep reading the local files only, since an ADR can refine only what its own branch holds.
- Two branches appending rows to `design/adrs/README.md` can still conflict textually at sync, as in S-0204's instance. With distinct numbers, resolving that conflict means keeping both rows, with no renumbering, so I left it out of scope.
- The forecast is 25m, up from flai's 18m by comparison with S-0252's 31m.

### 2026-10-05T05:46:11Z alex
Resolved.
