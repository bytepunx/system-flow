---
id: T-2012
type: task
nature: feature
title: A formula starts from a role whose file permits formula run with its name, which the guard allows for that formula alone, and flai formula run --ready makes its instance ready
status: backlog
parent: S-0422
owner: alex
created: 2026-10-21T12:53:13Z
updated: 2026-10-21T12:53:13Z
transitions: []
stream: S-0422
tags: [cli]
touches: [flai/cmd/formula.go, flai/cmd/formula_test.go, flai/internal/guard/roles.go, flai/internal/guard/roles_test.go, docs/users/flai-reference.md, docs/operators/settings.md]
---
# T-2012 A formula starts from a role whose file permits formula run with its name, which the guard allows for that formula alone, and flai formula run --ready makes its instance ready

## Work

- `flai formula run <name>` gains `--ready`: once the instance passes `flai check --strict`, it is moved to `ready` by the caller, so that `flai serve` starts its agent. Without the flag the instance stays in the backlog, as S-0421 makes it. A missing or undeclared input is still refused before anything is written.
- A role file may list `formula run <name>` among the commands it permits (S-0414's format, S-0416's guard). The guard's role path lets a defined role, the orchestrator under its role file among them, run `flai formula run <name>` for the formulas its file names, with `--ready` and `--input`, and refuses every other formula and every other `flai formula` write.
- Test the flag, and the guard with a role file that permits `formula run publish`: `flai formula run publish --ready` passes and `flai formula run import` is refused.
- Regenerate `docs/users/flai-reference.md` and the flag table in `docs/operators/settings.md` with `make flai-reference`.

Waits for nothing: it shares no path with the other tasks of the first layer. The schedule starter moves its instance to ready the same way.

## Done when

`go test ./cmd/... ./internal/guard/...` passes, and the settings-index test passes.

## Notes

Written by the planner on 2026-10-21, on TH-0471's answer 5: the schedule and role starters of T-1778 were split, this task taking the role starter and the flag that makes an instance ready.
