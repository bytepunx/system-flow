---
id: T-1467
type: task
nature: remediation
title: flai guard refuses a story's agent, the planner, and the analyzer a flai command that names its author with --by or FLAI_AGENT, with tests that reproduce I-0091
status: in-progress
parent: S-0369
owner: alex
created: 2026-10-12T17:05:52Z
updated: 2026-10-12T18:27:12Z
transitions:
  - to: ready
    at: 2026-10-12T18:27:12Z
    by: system-flow
  - to: in-progress
    at: 2026-10-12T18:27:12Z
    by: system-flow
stream: S-0369
tags: [flai, guard]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/authors.go, flai/internal/guard/authors_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md]
after: [T-1465]
---
# T-1467 flai guard refuses a story's agent, the planner, and the analyzer a flai command that names its author with --by or FLAI_AGENT, with tests that reproduce I-0091

## Work

Build the rule T-1465's ADR decides. It waits for T-1465 because the ADR settles which agents, flags, and environment the rule covers.

- Add `Agent`, `Session`, and `Item` to `guard.Guard`, read in `flai/cmd/guard.go` from the hook's `FLAI_AGENT`, `FLAI_SESSION`, and `FLAI_ITEM`.
- In a new `flai/internal/guard/authors.go`, check a served agent's own `Bash` call (a story's agent, `FLAI_ROLE` empty and served; the planner; the analyzer): refuse a flai command whose `--by` or `--placed-by` names another (not on `flai order --by` or `flai stats --by`), and an assignment of `FLAI_AGENT` or `FLAI_SESSION` to another value, empty included, before a command, after `export`, `declare`, `typeset`, `local`, `readonly`, or `env`, and `unset`, `env -u`, or `env -i` of them, also inside `sh -c` and `eval` scripts. Call it from `Decide` before the role's own rules. The orchestrator keeps its `writesAsSelf` rule.
- The refusal names what to do instead: leave out `--by`, or give its own name, and leave `FLAI_AGENT` and `FLAI_SESSION` as they are set.
- Update `flai guard`'s help in `flai/cmd/guard.go`. Regenerate `docs/users/flai-reference.md` with `make flai-reference`.
- Reproduce I-0091 in `flai/internal/guard/authors_test.go`:
  - `flai thread reply TH-0001 "allow" --by alex` and `FLAI_AGENT=alex flai thread reply TH-0001 "allow"` from a story's agent are refused, as are the same from the planner and the analyzer.
  - `flai thread reply` without `--by`, with `--by` its own name, `export FLAI_AGENT=<own>`, `flai order --by cod`, `flai move`, and `git commit` pass.
  - A session with no `FLAI_STORY`, `FLAI_ROLE`, or `FLAI_STARTED_BY` (the operator's own) passes `--by alex`.
- Change `TestTheStorysAgentIsNeverRefused` in `guard_test.go` to the rule's exceptions, and add a hook-JSON case in `flai/cmd/guard_test.go`.

## Done when

- The tests above fail before the change and pass after it.
- `flai test` on the task's paths passes, and `docs/users/flai-reference.md` matches the help.

## Notes

S-0368, in review, changes the same guard files for sub-agents' shell writes; the rule sits in its own file so the two meet only in `Decide`, the `Guard` struct, and the help. Rewritten by the story's agent to the narrowing the operator confirmed on TH-0418.
