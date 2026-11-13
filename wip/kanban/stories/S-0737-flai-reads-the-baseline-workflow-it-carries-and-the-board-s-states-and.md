---
id: S-0737
type: story
nature: feature
title: flai reads the baseline workflow it carries, and the board's states and transition map come from it with nothing else changed
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T15:58:44Z
updated: 2026-11-13T16:13:06Z
transitions: []
tags: [cli]
topics: [cli]
touches: [flai/internal/workflow/workflow.go, flai/internal/workflow/workflow_test.go, flai/internal/workflow/baseline.go, flai/internal/workflow/baseline_test.go, flai/internal/workflow/baseline/board.md, flai/internal/workitem/item.go, flai/internal/workitem/rules.go, flai/internal/workitem/transition.go, flai/internal/workitem/repo.go, flai/internal/workitem/rules_test.go, design/system/flai-cli.md]
after: [S-0736]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: sum
  seconds: 0
  models: []
  strategic:
    - kind: orchestrator
      seconds: 12
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 4
          output: 16
          cache_read: 108315
          cache_write: 18020
          cost: 0.0361
draft: true
forecast:
  duration: 33m
  delivery: 2026-11-13T18:25:00Z
  basis: "Its own forecast of 33m; 7th in the pull order with an in-progress limit of 3, behind S-0760, S-0759, S-0757, S-0241, S-0730 and S-0736."
  by: flai
  at: 2026-11-13T16:13:06Z
---
# S-0737 flai reads the baseline workflow it carries, and the board's states and transition map come from it with nothing else changed

## Goal

flai holds the board as data. A new package parses a workflow file in the format S-0736's ADR fixes, flai embeds its baseline file as it embeds the baseline role files, and `workitem` takes its list of states and its transition map from that baseline instead of `item.go`'s closed list and `rules.go`'s `allowed` map. Every later story of E-0034 reads its kinds, effects, rules, and timers from this package.

## Acceptance criteria

- [ ] A package `flai/internal/workflow` parses a workflow file into its states (name, kind, item types, WIP limit, order), its transitions (from, to, who may make it, reason, rules, effects), and its timers, and refuses a file that does not parse, naming the field and the line
- [ ] flai embeds its baseline workflow file, which gives today's six states, today's transitions, and the baseline timers as the ADR states them
- [ ] `workitem.States`, the check that a status is a state, `allowed`, `Follows`, and `allowedFor` read the baseline workflow; the six state constants stay as the baseline's names for the code later stories move
- [ ] A test compares the baseline's states and transition map with today's, pair by pair, for every item type, so a change to the baseline that alters the board fails
- [ ] A test opens a repository with a workflow other than the baseline, so later stories can test their packages against renamed and added states
- [ ] `flai/internal/workflow` imports nothing of `workitem`, and every existing test passes unchanged
- [ ] `design/system/flai-cli.md` § Internal structure names the package and what it holds

## Tasks

- T-3188 A workflow package parses a workflow file into its states, transitions, and timers, and refuses one that does not parse
- T-3195 flai embeds the baseline workflow file, and a test holds it to today's six states and transition map pair by pair
- T-3196 flai-cli.md names the workflow package and what it holds
- T-3201 workitem takes its states and transition map from the baseline workflow, and a test hook opens a repository with another workflow

## Notes

- The baseline role files are embedded the same way: `flai/internal/roles/baseline.go` with `//go:embed baseline/*.md`. The template does not ship the baseline workflow until the project's own file is read (S-0755); until then the embedded copy is the only one.
- The path `flai/internal/workflow/baseline/board.md` assumes the ADR names the file `board.md`. Where it names it otherwise, the story's agent follows the ADR and renames the touch.

### Planning

- Touches, all files:
  - `flai/internal/workflow/workflow.go`, `workflow_test.go`, `baseline.go`, `baseline_test.go`, `baseline/board.md`: layout. A new package beside `flai/internal/roles`, modelled on `roles.go` and `baseline.go`.
  - `flai/internal/workitem/item.go`, `rules.go`, `transition.go`, `repo.go`, `rules_test.go`: layout. The states are at `item.go` lines 29 to 41, the map at `rules.go` 15 to 21, and `Repo` carries the workflow it runs.
  - `design/system/flai-cli.md`: design (criterion 7). A shared path, so it holds nothing.
  - Left out from `flai touches suggest`: the user and operator guides and `flaiover-dashboard.md` (16% to 73%). Nothing a user sees changes in this story.
- `after`: S-0736. This story parses the format S-0736's ADR fixes, its fields and allowed values, and writes the baseline that ADR states.
- Tasks: T-3188 (the parser and its lookups, the mover check included); then T-3195 (the baseline) and T-3196 (`flai-cli.md`) together; then T-3201 (`workitem` reads it, and the test hook later stories use). T-3196 and T-3195 share no path.
- Forecast 33m, as `flai forecast` gives it for size 18. The change to `workitem` is a swap of a map's source behind the same functions, and the baseline test makes it checkable.
- Cost of delay: not set yet, pending TH-1097, as on S-0736.
