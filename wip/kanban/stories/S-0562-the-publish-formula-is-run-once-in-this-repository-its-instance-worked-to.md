---
id: S-0562
type: story
nature: improvement
title: The publish formula is run once in this repository, its instance worked to review and accepted, and the release or the evaluation that stopped it recorded
status: backlog
parent: E-0024
owner: alex
created: 2026-11-02T04:14:44Z
updated: 2026-11-02T04:14:44Z
transitions: []
tags: [template]
topics: [orchestration]
after: [S-0423]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
draft: true
---
# S-0562 The publish formula is run once in this repository, its instance worked to review and accepted, and the release or the evaluation that stopped it recorded

## Goal

E-0024's `publish` formula, shipped by S-0423, is run once in this repository before the epic is accepted. This shows that a publish is a story worked, measured, and archived: the formula's `in: main` steps, its `--require-met` evaluation, and the `push` gate in `flai release --pending` all run for real.

## Acceptance criteria

- [ ] The formula is run once in this repository: its instance story is worked to review and accepted, and its Notes name the release it published, or, when nothing was pending, the evaluation that stopped it.
- [ ] This story's Notes record the instance's ID, each step's exit, and the versions and tags published, or the evaluation that stopped the instance.
- [ ] What the run finds wrong is fixed in a task of this story, or recorded as an issue.

## Tasks

## Notes

- Moved out of S-0423 on TH-0803 (T-1783 and S-0423's criterion 4). The run needs `flai release --evaluate --require-met` and the `push` gate, which exist only once S-0423 is accepted, and the instance's acceptance is the operator's. This follows D-0422, where a run that needs the story's own change on the host moves to a follow-up story.
- `after: S-0423`: this story reads S-0423's output, the `publish` formula file and the flai changes its steps run.
- How to run it: with the `formula` and `push` host actions on, run `flai formula run publish --ready` in the main checkout, or let the orchestrator start it under `publish` when `release_evaluate` is met. `flai serve` starts the instance's agent, which runs the `in: main` steps from the main checkout. If those steps run the installed `flai` rather than `scripts/flai.sh`, the run waits until S-0423's flai is published too.
- Criterion 1 needs the operator's `push` host action on. When no accepted story is unreleased, the run shows the evaluation stopping the instance cleanly.
