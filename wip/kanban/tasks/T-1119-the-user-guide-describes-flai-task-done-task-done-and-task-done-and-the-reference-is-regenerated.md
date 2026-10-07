---
id: T-1119
type: task
nature: improvement
title: The user guide describes flai task done, task_done, and task.done, and the reference is regenerated
status: done
parent: S-0269
owner: alex
created: 2026-10-06T22:54:00Z
updated: 2026-10-07T03:05:28Z
transitions:
  - to: ready
    at: 2026-10-07T03:02:56Z
    by: agent-S-0269
  - to: in-progress
    at: 2026-10-07T03:02:57Z
    by: agent-S-0269
  - to: done
    at: 2026-10-07T03:05:28Z
    by: agent-S-0269
stream: S-0269
tags: [docs]
touches: [docs/users/flai.md, docs/users/flai-reference.md, docs/operators/settings.md]
after: [T-1102, T-1107, T-1113]
usage:
  source: log
  seconds: 151
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 58
      output: 375
      cache_read: 2213858
      cache_write: 89737
      cost: 1.1537
---
# T-1119 The user guide describes flai task done, task_done, and task.done, and the reference is regenerated

## Work

Document the operation for users.

- In `docs/users/flai.md`, where the guide walks through working a story, replace the per-task run of commands with `flai task done T-nnnn -m "<message>"`. Say:
  - the steps it runs, and in what order;
  - that it stops at the first failure, and what to do on a refused sync or a failed check;
  - that `--log` gives another log entry, and `--json` the structured answer;
  - that the same operation is `task_done` over MCP and `task.done` on the host channel.
- Regenerate `docs/users/flai-reference.md` with `make flai-reference` (`scripts/flai-reference.sh`). Never edit it by hand.

This task waits for T-1102, T-1107, and T-1113, so that it describes what was built. It is layer 5, the last.

## Done when

- [ ] `docs/users/flai.md` describes `flai task done`, `task_done`, and `task.done`.
- [ ] `docs/users/flai-reference.md` is regenerated and lists `flai task done` with its flags.
- [ ] The markdown lint and `flai check --strict` pass.

## Notes

Drafted by the planner. The design side, `design/system/flai-cli.md`, is T-1080's.
