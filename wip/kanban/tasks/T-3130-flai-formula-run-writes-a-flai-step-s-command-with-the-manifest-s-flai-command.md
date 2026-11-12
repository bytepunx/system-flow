---
id: T-3130
type: task
nature: feature
title: flai formula run writes a flai step's command with the manifest's flai.command when it names one
status: done
parent: S-0681
owner: alex
created: 2026-11-12T00:42:45Z
updated: 2026-11-12T00:46:47Z
transitions:
  - to: ready
    at: 2026-11-12T00:43:03Z
    by: agent-S-0681
  - to: in-progress
    at: 2026-11-12T00:43:03Z
    by: agent-S-0681
  - to: done
    at: 2026-11-12T00:46:47Z
    by: agent-S-0681
stream: S-0681
tags: []
touches: [flai/cmd/formula.go, flai/cmd/formula_test.go, docs/users/flai.md, design/system/project-manifest.md, design/system/strategic-agents.md]
usage:
  source: log
  seconds: 224
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 26
      output: 5246
      cache_read: 2036527
      cache_write: 64103
      cost: 1.0252
---
# T-3130 flai formula run writes a flai step's command with the manifest's flai.command when it names one

## Work

In `flai/cmd/formula.go`, `formulaTaskSections` takes the manifest's `flai.command` (ADR-0201) and, when it is set, writes a flai step's `## Work` and `## Done when` with that command in place of the bare `flai`, and the `flai formula stop` line of a step with `stop` too, so that the story's agent runs the flai the tree builds, which knows the flags the formula's steps use (I-0332). With no `flai.command`, the task says `flai` as before. Test it in `flai/cmd/formula_test.go` with a fixture manifest naming `flai.command`, for a step in the worktree and one `in: main` with `stop`. Say so in `docs/users/flai.md` (A story from a formula), `design/system/project-manifest.md` (the `flai.command` rule), and `design/system/strategic-agents.md` (The steps).

Waits for nothing: it is the first layer.

## Done when

A formula run in a project whose manifest names `flai.command` writes each flai step's task with that command, a test shows it and fails without the change, and the three documents say so.

## Notes
