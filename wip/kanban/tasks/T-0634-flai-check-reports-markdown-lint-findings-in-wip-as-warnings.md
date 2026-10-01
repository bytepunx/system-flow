---
id: T-0634
type: task
nature: remediation
title: flai check reports markdown lint findings in wip as warnings
status: ready
parent: S-0179
owner: arobson
created: 2026-10-01T08:23:14Z
updated: 2026-10-01T08:23:32Z
transitions:
  - to: ready
    at: 2026-10-01T08:23:32Z
    by: claude-opus-5-5
stream: S-0179
tags: []
touches: [flai/internal/check]
---
# T-0634 flai check reports markdown lint findings in wip as warnings

## Work

`flai check` lints every markdown file under the wip folder with `mdlint` when the project has a markdownlint configuration, and reports each finding as a warning with its rule and line, so `--strict` fails on it. `flai story new --body-stdin`, `flai edit`, and `item_edit` refuse a change that brings a new one, through the check they already run.

## Done when

- [ ] A wip file with a finding gives a warning naming the rule and line; a project without a markdownlint configuration gives none
- [ ] An edit whose body breaks a rule is refused with the rule and line
- [ ] `make test` passes

## Notes

Part of S-0179.
