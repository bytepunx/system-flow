---
id: T-1468
type: task
nature: improvement
title: docs/users/flai.md says who a thread entry's author is and that flai guard refuses a served agent --by and FLAI_AGENT on flai
status: done
parent: S-0369
owner: alex
created: 2026-10-12T17:06:01Z
updated: 2026-10-12T18:41:16Z
transitions:
  - to: ready
    at: 2026-10-12T18:40:51Z
    by: system-flow
  - to: in-progress
    at: 2026-10-12T18:40:51Z
    by: system-flow
  - to: done
    at: 2026-10-12T18:41:16Z
    by: agent-S-0369
stream: S-0369
tags: [flai, docs]
touches: [docs/users/flai.md]
after: [T-1466, T-1467]
usage:
  source: log
  seconds: 25
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 16
      output: 7105
      cache_read: 1096579
      cache_write: 34761
      cost: 0.5952
---
# T-1468 docs/users/flai.md says who a thread entry's author is and that flai guard refuses a served agent --by and FLAI_AGENT on flai

## Work

Update `docs/users/flai.md` in the same story that changes the commands, as `documentation.md` asks:

- Under `### Threads`, state the author rule: `--by`, else `FLAI_AGENT`, else the config's author, as T-1466 makes it.
- Say that a story's agent, the planner, and the analyzer cannot name an author: `flai guard` refuses them `--by` and an `FLAI_AGENT` set for flai, as T-1467 builds it. Their entries always carry their own name. Link the ADR from T-1465.
- In the guard's section, under `## Sub-agents` or wherever the guard's refusals are listed, add what it now refuses a story's own agent.
- Check `### Serving agents over MCP` for any sentence the change makes untrue.

It waits for T-1466 and T-1467 so that it describes what they built.

## Done when

- `docs/users/flai.md` describes the author rule and the guard's new refusal as built.
- It links the ADR and passes the markdown lint (`flai test docs/users/flai.md`).

## Notes

`docs/operators/settings.md` lists `FLAI_AGENT` under `### Read by flai`. Change it only if the guard now reads it: `flai task done` widens the touches then.
