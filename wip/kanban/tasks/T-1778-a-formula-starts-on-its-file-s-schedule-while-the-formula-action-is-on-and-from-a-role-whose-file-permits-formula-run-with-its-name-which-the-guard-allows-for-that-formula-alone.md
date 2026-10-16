---
id: T-1778
type: task
nature: feature
title: A formula starts on its file's schedule while the formula action is on, and from a role whose file permits formula run with its name, which the guard allows for that formula alone
status: backlog
parent: S-0422
owner: alex
created: 2026-10-16T05:48:01Z
updated: 2026-10-16T05:48:06Z
transitions: []
stream: S-0422
tags: [cli]
touches: [flai/internal/serve/formula_schedule.go, flai/internal/serve/formula_schedule_test.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, flai/internal/formulas/formulas.go, flai/internal/formulas/formulas_test.go]
after: [T-1774]
---
# T-1778 A formula starts on its file's schedule while the formula action is on, and from a role whose file permits formula run with its name, which the guard allows for that formula alone

## Work

- A formula whose file names a `schedule` is run by `flai serve` each time it comes round while the `formula` action is on, through the shared cadence, as `flai formula run` would run it with no inputs. A formula that needs inputs is refused a schedule by its check.
- A role file may list `formula run <name>` among its commands. The guard's role path allows that command for the named formula only.
- Test both starters.

Waits for the step runner. Runs beside the role step, failure, and acceptance tasks.

## Done when

`go test ./internal/serve/... ./internal/guard/...` passes.

## Notes
