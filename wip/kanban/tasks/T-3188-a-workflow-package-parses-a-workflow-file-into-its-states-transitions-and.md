---
id: T-3188
type: task
nature: feature
title: A workflow package parses a workflow file into its states, transitions, and timers, and refuses one that does not parse
status: backlog
parent: S-0737
owner: alex
created: 2026-11-12T16:05:19Z
updated: 2026-11-12T16:10:54Z
transitions: []
stream: S-0737
tags: [cli]
touches: [flai/internal/workflow/workflow.go, flai/internal/workflow/workflow_test.go]
---
# T-3188 A workflow package parses a workflow file into its states, transitions, and timers, and refuses one that does not parse

## Work

Add `flai/internal/workflow` with the types the ADR fixes: `Workflow`, `State` (name, kind, item types, WIP limit, order), `Transition` (from, to, who, reason, rules, effects), and `Timer`, and `Parse`, which reads the file's front matter as `roles.Parse` reads a role file and refuses a missing field, an unknown kind or effect, or a value outside the allowed ones, naming the field and the line. Give lookups the later stories need: a state by name, its kind, the transitions from a state for an item type, the transition with a given effect, and whether a given mover (the operator, a story's agent, a role, flai) may make a transition, with the words a refusal gives. Import nothing of `workitem`.

## Done when

`workflow_test.go` parses a valid file, refuses each kind of invalid one with the field and the line, and answers each lookup, the mover check included.

## Notes
