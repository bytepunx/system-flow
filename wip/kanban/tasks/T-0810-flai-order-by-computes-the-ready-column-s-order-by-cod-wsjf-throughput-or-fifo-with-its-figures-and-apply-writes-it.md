---
id: T-0810
type: task
nature: feature
title: flai order --by computes the ready column's order by cod, wsjf, throughput, or fifo with its figures, and --apply writes it
status: backlog
parent: S-0217
owner: alex
created: 2026-10-04T19:08:56Z
updated: 2026-10-04T19:08:56Z
transitions: []
stream: S-0217
tags: [flai]
touches: [flai/internal/workitem/policy.go, flai/internal/workitem/policy_test.go, flai/cmd/order.go, flai/cmd/order_test.go]
---
# T-0810 flai order --by computes the ready column's order by cod, wsjf, throughput, or fifo with its figures, and --apply writes it

## Work

Add a pure function to `flai/internal/workitem/policy.go`. It takes stories, their current board order, and a policy, and returns the stories in order, each with the figure it was ordered by:

- `cod`: cost of delay `value`, highest first.
- `wsjf`: `value` over the forecast `duration` in hours, highest first.
- `throughput`: forecast `duration`, shortest first.
- `fifo`: `created`, oldest first.

A story missing the figure its policy needs goes after the stories that have it, in its current board order, and names the missing figure. Ties keep the current board order. The accessors are in `planning.go`.

Give `flai order` (`flai/cmd/order.go`) `--by <policy>`, which:

- refuses to run with `--before`, `--after`, `--top`, or `--bottom`
- computes the ready column's order
- prints position, ID, figure, and title, or JSON with `--json`

With `--apply`, write `board.md`'s ready order through the functions dragging uses (`workitem/order.go`, `Place`). Without it, write nothing.

This task waits for nothing. The policy is given by the flag, and the manifest's setting is the first task's. It runs with that task, whose paths it does not share.

## Done when

- fixture tests pin each of the four orders, including a story missing its figure and a tie
- `flai order --by cod` prints the order with its figures and leaves `board.md` unchanged
- with `--apply`, `board.md`'s ready order is the computed one
- `--by` with a placement flag is refused
- `go test ./internal/workitem/ ./cmd/` passes

## Notes
