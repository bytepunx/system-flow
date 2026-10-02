---
id: S-0194
type: story
nature: feature
title: An experiment story is accepted with its results document under design/experiments
status: done
owner: arobson
created: 2026-10-02T09:46:12Z
updated: 2026-10-02T10:53:57Z
transitions:
  - to: ready
    at: 2026-10-02T10:06:38Z
    by: alex
  - to: in-progress
    at: 2026-10-02T10:07:23Z
    by: agent-S-0194
  - to: review
    at: 2026-10-02T10:32:11Z
    by: agent-S-0194
  - to: done
    at: 2026-10-02T10:53:57Z
    by: alex
tags: [flai, dashboard, template]
topics: [release]
touches: [flai/internal/preview, flai/internal/release, flai/cmd/accept.go, flai/cmd/accept_research_test.go, flai/internal/experiment, flai/internal/check, flaiover/src, template, design/experiments, CLAUDE.md, design/system/workflow.md, design/system/work-hierarchy.md, design/system/flai-cli.md, design/system/flaiover-dashboard.md, design/system/repository-layout.md, docs/users, docs/operators, design/README.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1516
  models:
    - model: claude-haiku-4-5-20251001
      input: 202
      output: 9807
      cache_read: 1337559
      cache_write: 83560
      cost: 0.2874
    - model: claude-opus-5-5
      input: 260
      output: 62431
      cache_read: 18064375
      cache_write: 208988
      cost: 6.5344
    - model: claude-sonnet-5-5
      input: 42
      output: 10824
      cache_read: 703086
      cache_write: 98936
      cost: 0.4963
---
# S-0194 An experiment story is accepted with its results document under design/experiments

## Goal

[ADR-0066](../../../design/adrs/0066-an-experiment-story-is-accepted-like-any-other-and-records-its-results-in-a.md) (2026-10-02) reverses ADR-0025's experiment clause: an `experiment` story is accepted like any other nature and records its results in a document under `design/experiments/`, and it contributes nothing to a component's release. Today acceptance refuses every experiment before merging (`flai/internal/preview/accept.go` ~91), the dashboard disables its accept button, and `design/experiments/` does not exist. The conventions already say the new rule (`git.md`, `documentation.md`). S-0176, an experiment in progress, needs this to be accepted.

## Acceptance criteria
- [x] `flai accept`, `flai move <story> done`, and the dashboard's acceptance accept an `experiment` story that has its results document, and refuse one without it, before anything is merged, naming the document expected
- [x] The results document is `design/experiments/<S-nnnn>-<slug>.md` with front matter (`title`, `updated`, `status`, the story's ID) and sections for the hypothesis, the success measure, what was done, the results, and the recommendation (adopt, adapt, or drop); `flai check` validates it, and a template for it ships with the template
- [x] Publishing gives no component a bump on an experiment's account, as for research; the release plan names it as landing unreleased
- [x] `design/experiments/README.md` exists here and in the template; the layout table in `CLAUDE.md` and the template's, and `design/system/repository-layout.md`, list the folder
- [x] `design/system/workflow.md`, `design/system/work-hierarchy.md`, `design/system/flai-cli.md`, `design/system/flaiover-dashboard.md`, and the user and operator guides describe accepting an experiment
- [x] Tests cover an experiment accepted with its document, refused without it, and left out of a release

## Tasks
- T-0683 Acceptance takes an experiment with its results document and refuses one without it, and the release plan gives it no bump
- T-0684 flai check validates design/experiments documents, and the folder, its README, and its document template ship here and in the template
- T-0685 The dashboard accepts an experiment and says when its results document is missing
- T-0686 The design and the user and operator guides describe accepting an experiment
- T-0687 The dashboard's ADR list test allows a proposed ADR in the monorepo

## Notes

Decided by the designer on 2026-10-02 while reviewing their convention edits.
