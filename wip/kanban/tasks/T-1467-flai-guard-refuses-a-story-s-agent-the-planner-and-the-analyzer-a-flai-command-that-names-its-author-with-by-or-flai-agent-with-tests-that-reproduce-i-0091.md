---
id: T-1467
type: task
nature: remediation
title: flai guard refuses a story's agent, the planner, and the analyzer a flai command that names its author with --by or FLAI_AGENT, with tests that reproduce I-0091
status: backlog
parent: S-0369
owner: alex
created: 2026-10-12T17:05:52Z
updated: 2026-10-12T17:05:52Z
transitions: []
stream: S-0369
tags: [flai, guard]
touches: [flai/internal/guard/guard.go, flai/internal/guard/guard_test.go, flai/internal/guard/thread_test.go, flai/cmd/guard.go, flai/cmd/guard_test.go, docs/users/flai-reference.md]
after: [T-1465]
---
# T-1467 flai guard refuses a story's agent, the planner, and the analyzer a flai command that names its author with --by or FLAI_AGENT, with tests that reproduce I-0091

## Work

Build the rule T-1465's ADR decides. It waits for T-1465 because the ADR settles which agents, commands, and environment the rule covers.

- In `Decide` in `flai/internal/guard/guard.go`, look at a story's own agent's `Bash` calls that run flai, which today pass untouched. Refuse those the ADR names: `--by` or `--placed-by` on a writing command, and `FLAI_AGENT` or `FLAI_SESSION` set for flai. Every other call of the story's agent still passes.
- Apply the same refusal to the planner's and the analyzer's rules, which today let `thread new` and `reply` pass with any `--by`. Keep the orchestrator's `writesAsSelf` rule as it is.
- Make the refusal name what to do instead: the MCP tool, or the command without `--by`, which writes the agent's own name.
- Update `flai guard`'s help in `flai/cmd/guard.go` to say what it refuses a story's own agent. Regenerate `docs/users/flai-reference.md` with `make flai-reference`. If the guard needs more of the run's identity than `FLAI_ROLE` and `FLAI_STORY`, read it there.
- Reproduce I-0091 in tests:
  - `flai thread reply TH-0001 "allow" --by alex` from a story's agent is refused.
  - `FLAI_AGENT=alex flai thread reply TH-0001 "allow"` from a story's agent is refused.
  - The same commands from the planner and the analyzer are refused.
  - `flai move`, `flai thread reply` without `--by`, and `git commit` from a story's agent pass.
  - A session with no `FLAI_STORY` and no `FLAI_ROLE` (the operator's own) passes `--by alex`.
- Change `TestTheStorysAgentIsNeverRefused` in `guard_test.go` to the rule's exceptions. Put the thread cases in `thread_test.go` and the hook-JSON cases in `flai/cmd/guard_test.go`.

## Done when

- The tests above fail before the change and pass after it.
- `flai test` on the task's paths passes, and `docs/users/flai-reference.md` matches the help.

## Notes

S-0368 changes the same guard files for sub-agents' shell writes. Whichever story runs second syncs onto the other.
