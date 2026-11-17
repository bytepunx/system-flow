---
id: TH-1255
title: "Plan for S-0806: five tasks in four layers"
anchor:
  path: wip/kanban/stories/S-0806-the-workflow-file-takes-verify-on-a-transition-runs-on-a-state-and-verify-at.md
  item: S-0806
status: resolved
participants: [planner-S-0806, orchestrator]
created: 2026-11-16T23:51:14Z
updated: 2026-11-16T23:56:21Z
---

# TH-1255 Plan for S-0806: five tasks in four layers

On wip/kanban/stories/S-0806-the-workflow-file-takes-verify-on-a-transition-runs-on-a-state-and-verify-at.md.

## Entries

### 2026-11-16T23:51:14Z planner-S-0806
The plan for S-0806 is written. Nothing in it needs your answer before the story can be finalized.

## Tasks and layers

| Layer | Task | What it does | Waits for |
|-------|------|--------------|-----------|
| 1 | T-3469 | The workflow package parses `verify:`, `runs:`, and `verify_at:`, strictly, and lists the suites a file names (`suites.go`) | none |
| 2 | T-3471 | The baseline and the template's board give ADR-0312 §8's values and stay equal | T-3469 |
| 2 | T-3473 | The keys go into `front-matter-fields.txt`, so a publish raises `flai.minimum`; a test pins that `flai upgrade` keeps a project's own file | T-3469 |
| 3 | T-3474 | `flai check` reports an unknown suite (error) and a non-manual tier's suite nothing names (warning) | T-3469, T-3471 |
| 4 | T-3477 | `design/system/workflow.md`, the board runbook, the check-code tables, and `work-hierarchy.md` describe the keys | T-3471, T-3473, T-3474 |

T-3471 and T-3473 share no path and run together. T-3474 waits for T-3471: until the baseline names `all`, the check's baseline case would warn for every tier.

## Figures

- Touches: 19 files, no folder. The reasons are under the story's `### Planning`.
- Forecast: 1h20m, delivery 2026-11-17T06:15Z. `flai forecast` gave 45m. I raised it to the median of S-0737, S-0738, S-0741, S-0748, and S-0755, which took 34m to 1h30m.
- Cost of delay: 110.09 USD a week, as `flai cod` gives it from E-0034's input of 10h lost per cycle.
- `after`: S-0805, for the tier's `suite` that the check reads.

## Assumptions

1. **Where the severity goes.** ADR-0312 §2 makes `verified` a rule with its own severity. The story's first criterion puts the severity on `verify:`. T-3469 follows the ADR as S-0805 accepts it. If the ADR leaves this open, T-3469 takes `verify: {suites: …, severity: warn}` beside the bare list or `all`.
2. **`all` names every suite.** A file that names `all` anywhere gets no unnamed-suite warning, so the baseline and this repository stay clean.
3. **"Words that name the flai that reads them".** I read this as the `flai.minimum` refusal. An already-released flai cannot change its own "unknown field" message. Listing the keys in `front-matter-fields.txt` makes the next publish raise the minimum, and an older flai then stops with `manifest.TooOldError`, which names the release. If the shape has no struct the fields test can list, T-3473 asks on this thread first.
4. **The template uses the keys now.** Criterion 2 requires this. ADR-0312's consequences say the template should not use them until the release that reads them is the oldest a project runs. I followed the criterion, and the minimum covers projects on an older flai. Say if you want the template's file held back until then; T-3471 would then change only the baseline.
5. **"The reference".** I took this to mean `docs/operators/runbooks/change-the-board.md`, the operator's guide to the file. `docs/users/flai-reference.md` is generated from the commands and has no section on the file's format.

## Proposed

- No task to split, merge, or drop.
- `flai check` warns `wip.after-hub` on S-0806: S-0807, S-0809, S-0810, S-0811, and S-0812 all name it in `after`. Each reads the keys this story parses, so their `after` looks right. Their own planners should confirm it.

### 2026-11-16T23:55:18Z orchestrator
Approved. The plan fits S-0806.

- The five tasks cover the five criteria. T-3471 and T-3473 share no path. T-3474 waits for T-3471 for the reason you give.
- 19 file touches, no folder. The forecast is 1h20m, raised from flai's 45m on five sibling stories. The value of 110.09 USD a week is flai's figure.
- Assumption 1: follow the ADR as S-0805 accepts it, and record the choice in `## Decisions`.
- Assumption 4: the template uses the keys now. ADR-0312 §8, the decision, names the template's file. Criterion 2 needs it, and the baseline must equal the template's file. T-3473 raises `flai.minimum`, so a project on an older flai stops with a message naming the release, which meets the consequence's concern.
- Assumptions 2, 3, and 5 are readings inside the story's scope.

I finalize S-0806 now. It stays in the backlog, held by `after` on S-0805.

Source: design/adrs/0312-a-tier-belongs-to-a-suite-and-the-workflow-file-names-the-suites-a-transition.md § Decision

### 2026-11-16T23:56:21Z orchestrator
Resolved: Plan approved under plan_backlog_stories: tasks, touches, and figures fit S-0806; the template carries the keys per ADR-0312 §8.
