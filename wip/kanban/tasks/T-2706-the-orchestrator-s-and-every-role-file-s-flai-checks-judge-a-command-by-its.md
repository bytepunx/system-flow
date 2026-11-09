---
id: T-2706
type: task
nature: remediation
title: The orchestrator's and every role file's flai checks judge a command by its arguments, its redirections and their targets left out
status: in-progress
parent: S-0613
owner: alex
created: 2026-11-09T08:11:25Z
updated: 2026-11-09T08:25:13Z
transitions:
  - to: ready
    at: 2026-11-09T08:25:13Z
    by: agent-S-0613
  - to: in-progress
    at: 2026-11-09T08:25:13Z
    by: agent-S-0613
stream: S-0613
tags: [flai, guard]
touches: [flai/internal/guard/guard.go, flai/internal/guard/roles.go, flai/internal/guard/orchestrate_test.go, flai/internal/guard/roles_test.go, flai/internal/guard/guard_test.go, flai/internal/guard/hostconfig.go, flai/internal/guard/hostconfig_test.go, flai/internal/guard/shared.go, flai/internal/guard/call.go]
---
# T-2706 The orchestrator's and every role file's flai checks judge a command by its arguments, its redirections and their targets left out

## Work

The cause: `commands` in `flai/internal/guard/guard.go` returns each simple command's `words`, which keep a redirection such as `2>&1` or `2>/dev/null` as a word. `Guard.refuse` hands those words to the orchestrator's rules and to a role file's (`asRole` in `flai/internal/guard/roles.go`), so `positionals` counts `2>&1` as a story in `flai order --by throughput --apply 2>&1`, and `subcommands` reads `2>/dev/null` as the subcommand of `flai board --json 2>/dev/null`. `parse` already splits each command's `args`, its words without its redirections and their targets (S-0368).

- Add a splitter beside `commands` that returns each simple command's `args`, and use it where the orchestrator's (`Guard.orchestrate`) and a role file's (`Guard.asRole`) shell checks call `commands`. Keep `commands` as it is for the callers and tests that read the words (S-0368).
- Leave the write check for read-only sub-agents as it is: it reads `targets`, and a redirection to a file is still a write there.
- Add tests to `flai/internal/guard/orchestrate_test.go` that reproduce I-0230's instances: `flai order --by throughput --apply 2>&1`, `flai order --by throughput --apply --json 2>&1 | jq .`, `cd <repo> && flai order --by throughput --apply 2>&1`, `flai board --json 2>/dev/null`, `flai queue --json 2>/dev/null`, and `flai serve status 2>&1` pass with the permissions that allow them; and that `flai order S-0001 --top 2>&1` is still refused as a placement by hand, so the redirection neither allows nor refuses anything.
- Add a test to `flai/internal/guard/roles_test.go` that a planner's allowed `flai` read with `2>&1` passes as it does without it.
- Added by the story's agent at review of the plan: `subcommands` skips only words that begin with `-`, so a redirection written before the subcommand, as in `flai 2>/dev/null serve enable push` or `flai 2>&1 accept S-0001 --by orchestrator`, is read as the command, and a check that refuses by command may let it through. Use the same splitter in every caller that hands `commands(...)` or `parse(...)` words to `Guard.refuse`: `orchestrate`, `plan`, `analyze`, `curate`, and `operatorsOnly` in `guard.go`, `asRole` in `roles.go`, the host-configuration check in `hostconfig.go`, the shared-paths check in `shared.go`, the call in `call.go`, and `Guard.script`, so a shell's `-c` script is judged the same way. Add tests that a write with a redirection before its subcommand is refused as it is without one: for the orchestrator in `orchestrate_test.go` and for the host configuration in `hostconfig_test.go`.

It waits for no task.

## Done when

- [ ] The tests above fail on the code before the change and pass after it
- [ ] `flai test flai/internal/guard` passes

## Notes
