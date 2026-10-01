---
id: T-0633
type: task
nature: remediation
title: flai has a markdown linter for what it writes, configured by the project's markdownlint file
status: done
parent: S-0179
owner: arobson
created: 2026-10-01T08:23:14Z
updated: 2026-10-01T08:40:19Z
transitions:
  - to: ready
    at: 2026-10-01T08:23:32Z
    by: claude-opus-5-5
  - to: in-progress
    at: 2026-10-01T08:24:02Z
    by: claude-opus-5-5
  - to: done
    at: 2026-10-01T08:40:19Z
    by: claude-opus-5-5
stream: S-0179
tags: []
touches: [flai/internal/mdlint, scripts/mdlint-fixtures.sh, scripts/README.md, Makefile]
usage:
  source: log
  seconds: 977
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 71
      output: 44538
      cache_read: 7775461
      cache_write: 84411
      cost: 3.1214
---
# T-0633 flai has a markdown linter for what it writes, configured by the project's markdownlint file

## Work

A package `flai/internal/mdlint` loads the project's `.markdownlint.yaml` (or `.yml`, `.json`, `.jsonc`) and checks a document against the subset of markdownlint rules that what flai writes can break, with the options the file sets and inline `markdownlint-disable` comments honoured. No project configuration means no lint.

## Done when

- [x] Fixtures with findings give the rule and line markdownlint-cli2 gives for every rule the package implements
- [x] Every markdown file in this repository, which markdownlint-cli2 passes, gives no finding
- [x] `make test` and the lint pass

## Notes

Part of S-0179.
