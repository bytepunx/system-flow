---
id: T-0884
type: task
nature: feature
title: flai order records who placed a story by hand and when, and flai order --by --apply keeps a hand placement of the last day
status: backlog
parent: S-0219
owner: alex
created: 2026-10-05T04:46:19Z
updated: 2026-10-05T04:46:19Z
transitions: []
stream: S-0219
tags: [flai]
touches: [flai/internal/workitem/board.go, flai/internal/workitem/board_placed_test.go, flai/internal/workitem/policy.go, flai/internal/workitem/policy_test.go, flai/cmd/order.go, flai/cmd/order_test.go, design/adrs]
---
# T-0884 flai order records who placed a story by hand and when, and flai order --by --apply keeps a hand placement of the last day

## Work

Today `board.md` keeps the pull order as a list and nothing says who placed a story or when (`workitem/board.go`; `Reordered` in `changes.go` only compares two lists). The criterion that the orchestrator does not reorder a story the operator ordered by hand in the last day needs that record.

- Record each placement `flai order <story> --before|--after|--top|--bottom` makes: the story, who placed it (`--by`, else `FLAI_AGENT`, else the config author, as `flai edit` stamps), and when. The dashboard's drag calls `flai order`, so it is recorded too. Keep the record where every host and the metrics can read it, such as a `placed` map in `board.md`'s front matter beside `order`, and drop a story's entry when `flai move` takes it out of `ready` and `backlog`. Write an ADR for where it lives and why.
- Give S-0217's `flai order --by <policy> --apply` (`workitem/policy.go`, `cmd/order.go`) a window, such as `--keep-placed 24h`, defaulting to a day: a ready story placed within the window by anyone but an orchestrator keeps its place, and the rest of the column is ordered by the policy around it. The printed order marks each kept story and who placed it when.

This task waits for no other task of this story; it needs S-0217's `flai order --by`, which S-0219 waits for through S-0218. It runs with the other first-layer tasks, whose paths it does not share.

## Done when

- `flai order` records who placed the story and when, and `flai move` out of `ready` or `backlog` drops its record
- `flai order --by cod --apply` on a fixture board keeps a story placed by hand an hour ago where it is and reorders the others; one placed two days ago, or by an orchestrator, is reordered
- the ADR is accepted in `design/adrs` and listed in its README
- `go test ./internal/workitem/ ./cmd/` passes

## Notes
