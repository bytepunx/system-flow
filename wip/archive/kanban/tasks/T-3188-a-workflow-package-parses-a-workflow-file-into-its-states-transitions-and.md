---
id: T-3188
type: task
nature: feature
title: A workflow package parses a workflow file into its states, transitions, and timers, and refuses one that does not parse
status: done
parent: S-0737
owner: alex
created: 2026-11-12T16:05:19Z
updated: 2026-11-14T23:21:25Z
transitions:
  - to: ready
    at: 2026-11-14T23:15:52Z
    by: system-flow
  - to: in-progress
    at: 2026-11-14T23:15:53Z
    by: system-flow
  - to: done
    at: 2026-11-14T23:21:25Z
    by: agent-S-0737
stream: S-0737
tags: [cli]
touches: [flai/internal/workflow/workflow.go, flai/internal/workflow/workflow_test.go]
usage:
  source: log
  seconds: 332
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 30
      output: 14959
      cache_read: 3641915
      cache_write: 83801
      cost: 1.6441
---
# T-3188 A workflow package parses a workflow file into its states, transitions, and timers, and refuses one that does not parse

## Work

Add `flai/internal/workflow` with the types the ADR fixes: `Workflow`, `State` (name, kind, item types, WIP limit, order), `Transition` (from, to, who, reason, rules, effects), and `Timer`, and `Parse`, which reads the file's front matter as `roles.Parse` reads a role file and refuses a missing field, an unknown kind or effect, or a value outside the allowed ones, naming the field and the line. Give lookups the later stories need: a state by name, its kind, the transitions from a state for an item type, the transition with a given effect, and whether a given mover (the operator, a story's agent, a role, flai) may make a transition, with the words a refusal gives. Import nothing of `workitem`.

## Done when

`workflow_test.go` parses a valid file, refuses each kind of invalid one with the field and the line, and answers each lookup, the mover check included.

## Notes
