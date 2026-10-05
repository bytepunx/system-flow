---
id: T-0862
type: task
nature: remediation
title: item_edit and the dashboard's writes accept a touch that starts with a dot, by one touch rule in workitem
status: cancelled
parent: S-0260
owner: alex
created: 2026-10-05T03:13:41Z
updated: 2026-10-05T03:21:09Z
transitions:
  - to: cancelled
    at: 2026-10-05T03:21:09Z
    by: agent-S-0260
stream: S-0260
tags: [flai]
touches: [flai/internal/workitem/create.go, flai/internal/workitem/create_test.go, flai/internal/itemedit/itemedit.go, flai/internal/hostapi/writes.go, flai/internal/hostapi/writes_test.go, flai/internal/mcpserver/items_write_test.go, flai/cmd/edit_test.go]
---
# T-0862 item_edit and the dashboard's writes accept a touch that starts with a dot, by one touch rule in workitem

## Work

The cause of I-0071: `itemedit.cleanList` checks touches with `listValue`, `^[A-Za-z0-9][A-Za-z0-9 _./@+-]*$` (`flai/internal/itemedit/itemedit.go`), so a touch must start with a letter or digit, though the refusal says only "not starting with a dash". The dashboard's writes copy the rule for touches on create and edit (`listValue` in `flai/internal/hostapi/writes.go`). A path such as `.claude/agents/planner.md` or `.github/workflows` is refused by both, while `flai touches` and `item_new` take it.

- Add `CleanTouches` to `flai/internal/workitem/create.go`, beside `CleanTopics`. It trims each touch, drops a trailing slash, empties, and repeats, and accepts a first character that is a letter, a digit, or a dot. It refuses a touch that starts with a dash (a flag), one that is absolute, one with a `..` segment, and one with a comma or a character outside `_ . / @ + -`.
- Check touches with it in `itemedit` (the `ch.Touches` branch), so `flai edit --touches` and MCP `item_edit` share it. Tags keep `cleanList` and its rule.
- Check touches with it in the dashboard's create and edit writes, in place of `listValue`. Tags keep `listValue`.
- Test the rule in `create_test.go`. Reproduce I-0071 in `flai/internal/mcpserver/items_write_test.go`: `item_edit` setting `.claude/agents/planner.md` is accepted and written. Add cases to `flai/cmd/edit_test.go` and `flai/internal/hostapi/writes_test.go`: a dot touch is accepted, and `--touches=--owner=eve` and `../x` are refused.

It waits for no task: the other tasks build on the rule it adds.

## Done when

- `item_edit`, `flai edit --touches`, and the dashboard's create and edit accept `.claude/agents/planner.md` as a touch.
- A touch that starts with a dash, is absolute, or has a `..` segment is refused by each of them, naming the touch.
- The new tests fail without the change and pass with it, and `scripts/flai-test.sh` passes.

## Notes
- 2026-10-05T03:21:09Z: moved to cancelled: duplicate of T-0865 and T-0866, which agent-S-0260 wrote on pulling S-0260 (TH-0124)
