---
id: T-1123
type: task
nature: improvement
title: flai plan's help, the MCP plan tool's description, and the conventions design say an epic's planner drafts its stories' tasks
status: done
parent: S-0300
owner: alex
created: 2026-10-06T23:07:16Z
updated: 2026-10-06T23:28:14Z
transitions:
  - to: ready
    at: 2026-10-06T23:27:12Z
    by: agent-S-0300
  - to: in-progress
    at: 2026-10-06T23:27:12Z
    by: agent-S-0300
  - to: done
    at: 2026-10-06T23:28:14Z
    by: agent-S-0300
stream: S-0300
tags: [planner, cli, docs]
touches: [flai/cmd/plan.go, flai/internal/mcpserver/plan.go, docs/users/flai-reference.md, design/system/conventions.md]
after: [T-1047]
usage:
  source: log
  seconds: 62
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 6
      output: 23
      cache_read: 462852
      cache_write: 3090
      cost: 0.2053
---
# T-1123 flai plan's help, the MCP plan tool's description, and the conventions design say an epic's planner drafts its stories' tasks

## Work

Criterion 3 changes what a planner run on an epic writes. Three texts outside the other tasks still say that it drafts only the stories:

- `flai/cmd/plan.go`, the `Long` help of `flai plan`: "For an epic with no stories it drafts the stories that deliver its outcome". Say that it drafts the stories and each one's tasks. For an epic with stories, say that it also drafts the tasks of each story it adds.
- `flai/internal/mcpserver/plan.go`, `planDescription`: "it drafts an epic's stories". Say the same in the tool's words.
- `design/system/conventions.md`, the `plan` row of the roles table: "drafts an epic's stories or enriches a story". Add the tasks.

Then regenerate `docs/users/flai-reference.md` with `scripts/flai-reference.sh`. Do not edit it by hand.

Use the same words as T-1047's design and convention text, so this task waits for T-1047 and is in layer 2. It shares no path with T-1050 or T-1051.

S-0261 is in progress and claims `flai/internal/mcpserver` and `docs/users/flai-reference.md`. Only the one constant in `plan.go` changes here. If `flai stream sync` stops on the reference, regenerate it.

## Done when

- [ ] `flai plan --help`, the MCP tool `plan`'s description, and the `plan` row in `design/system/conventions.md` say that planning an epic drafts its stories and their tasks.
- [ ] `docs/users/flai-reference.md` is regenerated and matches the help.
- [ ] `scripts/flai-test.sh`, `flai check --strict`, and the markdown lint pass.

## Notes

Added by planner-S-0300 when it revisited the plan on 2026-10-06. `grep` found the sentence in these files, and no other task names them.
