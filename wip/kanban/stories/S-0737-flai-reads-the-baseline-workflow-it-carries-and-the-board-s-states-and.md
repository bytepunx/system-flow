---
id: S-0737
type: story
nature: feature
title: flai reads the baseline workflow it carries, and the board's states and transition map come from it with nothing else changed
status: backlog
parent: E-0034
owner: alex
created: 2026-11-12T15:58:44Z
updated: 2026-11-12T15:58:44Z
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
draft: true
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

## Notes

- The baseline role files are embedded the same way: `flai/internal/roles/baseline.go` with `//go:embed baseline/*.md`. The template does not ship the baseline workflow until the project's own file is read (S-0750's work); until then the embedded copy is the only one.
- The path `flai/internal/workflow/baseline/board.md` assumes the ADR names the file `board.md`. Where it names it otherwise, the story's agent follows the ADR and renames the touch.
