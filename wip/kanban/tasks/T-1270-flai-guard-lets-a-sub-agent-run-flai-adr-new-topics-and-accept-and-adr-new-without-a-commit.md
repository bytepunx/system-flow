---
id: T-1270
type: task
nature: remediation
title: flai guard lets a sub-agent run flai adr new, topics, and accept, and adr_new, without a commit
status: in-progress
parent: S-0287
owner: alex
created: 2026-10-07T23:19:09Z
updated: 2026-10-08T07:14:08Z
transitions:
  - to: ready
    at: 2026-10-08T07:14:07Z
    by: agent-S-0287
  - to: in-progress
    at: 2026-10-08T07:14:08Z
    by: agent-S-0287
stream: S-0287
tags: [cli, flai, guard]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/cmd/guard_test.go]
---
# T-1270 flai guard lets a sub-agent run flai adr new, topics, and accept, and adr_new, without a commit

## Work

In `flai/internal/guard/guard.go`, let a sub-agent's rules (`subAgent`, through `reads` and its tables) pass `flai adr new` (with `--print-body` too), `flai adr topics`, and `flai adr accept` when none names `--commit` or `--autocommit`, and pass the MCP tool `adr_new` when its `commit` is not true. An ADR is a document, which ADR-0060 does not guard; committing it writes history and the story's touches, which stay the story's agent's. Leave the planner's, the orchestrator's, and the analyzer's rules as they are.

In `flai/internal/guard/guard_test.go`, cover each form for a sub-agent: allowed bare and with `--print-body`, refused with `--commit`, `--commit=true`, and `--autocommit`, and `adr_new` allowed without `commit` and refused with it. In `flai/cmd/guard_test.go`, reproduce I-0062 through the binary: a sub-agent's `flai adr new` and then `flai adr topics` on the same ADR both pass the guard.

It waits for no task: the rule is the fix, and the ADR and the docs describe it.

## Done when

- A sub-agent's `flai adr new`, `flai adr topics`, and `flai adr accept` without `--commit` or `--autocommit`, and its `adr_new` without `commit`, pass `flai guard`; each with a commit is refused with the reason.
- The tests above pass with `flai test` on the three files.

## Notes
