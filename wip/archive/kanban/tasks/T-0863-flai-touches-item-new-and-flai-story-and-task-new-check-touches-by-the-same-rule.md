---
id: T-0863
type: task
nature: remediation
title: flai touches, item_new, and flai story and task new check touches by the same rule
status: cancelled
parent: S-0260
owner: alex
created: 2026-10-05T03:13:50Z
updated: 2026-10-05T03:21:09Z
transitions:
  - to: cancelled
    at: 2026-10-05T03:21:09Z
    by: agent-S-0260
stream: S-0260
tags: [flai]
touches: [flai/internal/workitem/create.go, flai/internal/workitem/create_test.go, flai/cmd/touches.go, flai/cmd/touches_test.go, flai/internal/mcpserver/items_write_test.go]
after: [T-0862]
---
# T-0863 flai touches, item_new, and flai story and task new check touches by the same rule

## Work

I-0071 is a disagreement: `flai touches` (`flai/cmd/touches.go`) and `workitem.Create` (`opt.Touches`, behind `item_new` and `flai story new` and `flai task new --touches`) take any touch, while `item_edit` refused one. Once T-0862 gives `item_edit` the rule, make every way of setting touches apply it, so that no path one accepts another refuses.

- In `workitem.Create`, clean `opt.Touches` with `CleanTouches` and return its refusal, as it does `opt.Topics` with `CleanTopics`.
- In `flai touches <id> <paths>`, clean the paths with `CleanTouches` before saving, and refuse what it refuses, writing nothing.
- Test both in `flai/internal/workitem/create_test.go`, `flai/cmd/touches_test.go`, and `flai/internal/mcpserver/items_write_test.go`: a touch that starts with a dot is kept, and one that starts with a dash or has a `..` segment is refused.

It waits for T-0862, which adds `CleanTouches` and changes `create.go` and `items_write_test.go`, the files this task changes too.

## Done when

- `flai touches`, `item_new`, and `flai story new` and `flai task new --touches` accept `.claude/agents/planner.md` and refuse `--owner=eve` and `../x` with the same message as `item_edit`.
- `scripts/flai-test.sh` passes.

## Notes
- 2026-10-05T03:21:09Z: moved to cancelled: duplicate of T-0865 and T-0866, which agent-S-0260 wrote on pulling S-0260 (TH-0124)
