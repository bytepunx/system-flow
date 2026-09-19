---
id: T-0216
type: task
nature: feature
title: "release: research plans no release and says why, naming any component files landing without one; experiment still refuses"
status: done
parent: S-0053
owner: alex
created: 2026-09-19T07:15:37Z
updated: 2026-09-19T07:20:22Z
transitions:
  - to: ready
    at: 2026-09-19T07:17:00Z
    by: system-flow
  - to: in-progress
    at: 2026-09-19T07:17:00Z
    by: system-flow
  - to: done
    at: 2026-09-19T07:20:22Z
    by: system-flow
stream: S-0053
tags: []
---

# T-0216 release: research plans no release and says why, naming any component files landing without one; experiment still refuses

## Work
In `flai/internal/release`: `LevelFor` stops returning an error for `research` and the plan for a research story is an explicit no-release plan with the reason, not an error and not an empty plan that looks like nothing was touched. When the story's commits touched a component's files the plan lists those components as landing without a release, so the operator sees code reaching main unreleased. `experiment` keeps its refusal, reworded to say it alone stays on a branch. The printed plan, `--json`, and the dry run all carry it. Unit tests for `LevelFor` and `Compute`: research with no component files, research with component files, experiment.

## Done when
- A research story's plan says no release and why, and names components whose files it touched
- No tag, version, or changelog change is planned for research
- Experiment is still refused, with a message that names only experiment

## Notes
