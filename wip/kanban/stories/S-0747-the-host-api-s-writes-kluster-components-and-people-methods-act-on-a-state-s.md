---
id: S-0747
type: story
nature: feature
title: The host API's writes, kluster, components, and people methods act on a state's kind and who may make a transition instead of a state's name
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T16:00:29Z
updated: 2026-11-12T16:00:29Z
transitions: []
tags: [cli, dashboard]
topics: [cli, dashboard]
touches: [flai/internal/hostapi/writes_describe.go, flai/internal/hostapi/writes_methods.go, flai/internal/hostapi/writes_items.go, flai/internal/hostapi/writes_outcome.go, flai/internal/hostapi/writes_run.go, flai/internal/hostapi/writes_test.go, flai/internal/hostapi/kluster.go, flai/internal/hostapi/kluster_forward.go, flai/internal/hostapi/kluster_test.go, flai/internal/hostapi/components.go, flai/internal/hostapi/people.go, flai/internal/hostapi/people_test.go, flai/internal/hostapi/import.go, flai/internal/hostapi/hostapi.go, flai/internal/hostapi/settings.go]
after: [S-0737]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0747 The host API's writes, kluster, components, and people methods act on a state's kind and who may make a transition instead of a state's name

## Goal

What the dashboard asks of flai on the host follows the workflow. A move from the board is refused or allowed by who may make its transition, an acceptance from the board makes the acceptance transition, and the descriptions of each write, the Deployments page's namespaces, the components' worktrees, and the people's open work read a state's kind, so the dashboard's writes work on a workflow they have never seen.

## Acceptance criteria

- [ ] The host API's files in `flai/internal/hostapi` name no state: each reads a state's kind, a transition's effect, or who may make it from the workflow
- [ ] `item.move` from the dashboard refuses a transition the operator may not make there, naming the transition and who may make it, and makes the acceptance transition through acceptance as today
- [ ] `writes_describe.go` describes each move in the workflow's names and the reason it asks for, from the transition, not from a list of states
- [ ] Tests run the move, the acceptance, the descriptions, and the kluster and people reads against a workflow with a state renamed and one added between `in-progress` and `review`
- [ ] With the baseline, every existing `hostapi` test and flaiover's tests of the built flai pass unchanged

## Tasks

## Notes

- `hostapi` holds about 73 references in 12 files, 31 of them in `writes_describe.go`. The method that gives the dashboard the workflow's columns is the board's story's; this one changes no answer's shape.
