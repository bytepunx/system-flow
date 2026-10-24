---
id: T-2103
type: task
nature: feature
title: task_done's check lists at most ten notes outside the story and counts the rest in notes_omitted, with a test, the help, and the docs
status: done
parent: S-0492
owner: alex
created: 2026-10-23T23:48:10Z
updated: 2026-10-23T23:58:52Z
transitions:
  - to: ready
    at: 2026-10-23T23:49:37Z
    by: agent-S-0492
  - to: in-progress
    at: 2026-10-23T23:49:39Z
    by: agent-S-0492
  - to: done
    at: 2026-10-23T23:58:52Z
    by: agent-S-0492
stream: S-0492
tags: []
touches: [flai/internal/taskdone/taskdone.go, flai/internal/taskdone/taskdone_test.go, flai/cmd/task_done.go, flai/internal/mcpserver/task.go, design/system/flai-cli.md, docs/users/flai.md, docs/users/flai-reference.md]
after: [T-2102]
usage:
  source: log
  seconds: 553
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 32
      output: 6816
      cache_read: 2072329
      cache_write: 46603
      cost: 0.9237
---

# T-2103 task_done's check lists at most ten notes outside the story and counts the rest in notes_omitted, with a test, the help, and the docs

## Work

In `flai/internal/taskdone`, the check step keeps the first ten notes and counts the rest in `Check.NotesOmitted` (`notes_omitted`, omitted when zero). `flai task done` prints the whole count and says how many it left out and that `flai check --story` lists them. A test reproduces I-0202: a check with hundreds of notes outside the story answers at most ten, counts the rest, and the result fits `ctxpack.PartLimit`. Update the command's help, `docs/users/flai-reference.md` with `make flai-reference`, the MCP tool's description, `design/system/flai-cli.md`, and `docs/users/flai.md`. Waits for T-2102, whose ADR the docs link.

## Done when

The test fails on main and passes here; `flai test` passes on the paths changed; the help, the reference, the tool description, and the docs say what the answer carries.

## Notes
