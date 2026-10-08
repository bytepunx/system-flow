---
id: S-0345
type: story
nature: remediation
title: A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix
status: backlog
owner: alex
created: 2026-10-08T08:08:19Z
updated: 2026-10-08T08:40:10Z
transitions: []
tags: [flai]
topics: [testing]
touches: [flai/internal/mdlint/mdlint_test.go, flai/internal/mdlint/repo_test.go, scripts/lint-md.sh, scripts/README.md, design/system/flai-cli.md, design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md, design/issues/summary.md]
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
      seconds: 7
      estimated: true
      models:
        - model: claude-opus-5-5
          input: 5
          output: 39
          cache_read: 1177495
          cache_write: 6618
          cost: 0.2918
cost_of_delay:
  inputs:
    time_lost_per_cycle: 1h12m
    by: flai
    at: 2026-10-08T08:08:19Z
finalized:
  by: alex
  at: 2026-10-08T08:40:10Z
---
# S-0345 A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix

## Goal

This story remediates [I-0117](../../../design/issues/I-0117-a-story-s-integration-tier-lints-main-s-committed-wip-so-markdown-any-agent-commits-there-fails-every-story-s-close-out-until-main-commits-a-fix.md), "A story's integration tier lints main's committed wip, so markdown any agent commits there fails every story's close-out until main commits a fix". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [ ] The cause I-0117 describes no longer occurs, with a test that reproduces it where one fits
- [ ] I-0117 is closed with `flai issue close I-0117 --reason` saying what fixed it

## Tasks
- T-1363 TestRepositoryLintsClean scopes itself to the story in a close-out, so main's committed wip is a note, not a failure
- T-1364 lint-md.sh's whole-repository run leaves main's wip out in a close-out and lints only the wip files the story changes
- T-1365 Document that the repository lint test and lint-md.sh scope themselves to the story in a close-out
- T-1366 Close I-0117 with the reason that names the scoped lint test and lint-md.sh

## Notes

Cost of delay inputs set by flai from I-0117. time_lost_per_cycle 1h12m: 18m per occurrence × 4 occurrences ÷ 1 cycle of 168h (first reported 2026-10-08T04:12:50Z, 0.2 days before this story; under one cycle counts as one).
