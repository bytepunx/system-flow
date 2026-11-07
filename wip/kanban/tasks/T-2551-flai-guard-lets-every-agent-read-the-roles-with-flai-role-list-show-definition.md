---
id: T-2551
type: task
nature: improvement
title: flai guard lets every agent read the roles with flai role list, show, definition, and role alone, and keeps run and --write refused
status: done
parent: S-0568
owner: alex
created: 2026-11-07T08:18:56Z
updated: 2026-11-07T08:26:37Z
transitions:
  - to: ready
    at: 2026-11-07T08:20:56Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T08:20:57Z
    by: system-flow
  - to: done
    at: 2026-11-07T08:26:37Z
    by: agent-S-0568
stream: S-0568
tags: [flai, guard]
touches: [flai/internal/guard/guard.go, flai/internal/guard/roles_test.go, flai/cmd/guard_test.go]
usage:
  source: log
  seconds: 340
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 44
      output: 331
      cache_read: 4516419
      cache_write: 30347
      cost: 1.6278
---
# T-2551 flai guard lets every agent read the roles with flai role list, show, definition, and role alone, and keeps run and --write refused

## Work

Layer 1: this task waits for nothing.

The cause, from I-0245's one instance: `reads` in `flai/internal/guard/guard.go` finds no `role` entry in `cliReads` or `flagReads`, so `subAgent` refuses `scripts/flai.sh role definition planner` with the reason for commands that change work items. S-0553's task sub-agent met it on T-2351.

The proposed fix: add `role` to `flagReads`, with `--write` in `without` and a `form` that allows only the subcommands `list`, `show`, and `definition`, and no subcommand, which prints help. `flagReads` is the table for this: a command that reads in some forms and writes in others, as `verify` is. `cliReads` alone cannot refuse `flai role definition --write`, because it looks only at the subcommand.

- `flai role run` stays refused, because it starts an agent.
- `flai role definition --write` stays refused, in any form, because it writes files.
- `reads` is shared, so the planner, the orchestrator, the analyzer, the curator, and a defined role gain the same reads for their own calls, as they did for `flai formula list`.

If the code shows a simpler form that keeps these refusals, use it and say why in the narrative's `## Decisions`.

## Done when

- A test in `flai/internal/guard/roles_test.go` passes `flai role`, `flai role list --json`, `flai role show planner`, `flai role definition planner`, and `scripts/flai.sh role definition planner` for the explorer, the verifier, and a task sub-agent, and fails on the guard without this change.
- The same test refuses `flai role run reviewer`, `flai role definition --write`, and `flai role definition --write planner` to each of those sub-agents.
- The planner's, the orchestrator's, the analyzer's, and the curator's own calls pass the reads, and `flai role definition --write` is refused to each.
- A case in `flai/cmd/guard_test.go` feeds `flai guard` a hook input with an `agent_id`: `scripts/flai.sh role definition planner` exits 0 and `flai role run reviewer` exits 2.
- `flai test flai/internal/guard flai/cmd` passes.

## Notes
