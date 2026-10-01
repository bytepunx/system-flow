---
id: T-0634
type: task
nature: remediation
title: flai check reports markdown lint findings in wip as warnings
status: done
parent: S-0179
owner: arobson
created: 2026-10-01T08:23:14Z
updated: 2026-10-01T08:43:02Z
transitions:
  - to: ready
    at: 2026-10-01T08:23:32Z
    by: claude-opus-5-5
  - to: in-progress
    at: 2026-10-01T08:40:19Z
    by: claude-opus-5-5
  - to: done
    at: 2026-10-01T08:43:02Z
    by: claude-opus-5-5
stream: S-0179
tags: []
touches: [flai/internal/check]
usage:
  source: log
  seconds: 163
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 33
      output: 20771
      cache_read: 3626311
      cache_write: 39368
      cost: 1.4558
---
# T-0634 flai check reports markdown lint findings in wip as warnings

## Work

`flai check` lints every markdown file under the wip folder with `mdlint` when the project has a markdownlint configuration, and reports each finding as a warning with its rule and line, so `--strict` fails on it. `flai story new --body-stdin`, `flai edit`, and `item_edit` refuse a change that brings a new one, through the check they already run.

## Done when

- [x] A wip file with a finding gives a warning naming the rule and line; a project without a markdownlint configuration gives none
- [x] An edit whose body breaks a rule is refused with the rule and line
- [x] `make test` passes

## Notes

Part of S-0179.
