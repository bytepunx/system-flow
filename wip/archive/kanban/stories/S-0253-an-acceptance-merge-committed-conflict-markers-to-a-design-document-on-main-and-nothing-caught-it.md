---
id: S-0253
type: story
nature: remediation
title: An acceptance merge committed conflict markers to a design document on main, and nothing caught it
status: done
owner: alex
created: 2026-10-03T20:06:56Z
updated: 2026-10-05T04:04:28Z
transitions:
  - to: ready
    at: 2026-10-04T21:42:57Z
    by: alex
  - to: in-progress
    at: 2026-10-05T03:10:48Z
    by: agent-S-0253
  - to: review
    at: 2026-10-05T04:01:37Z
    by: agent-S-0253
  - to: done
    at: 2026-10-05T04:04:28Z
    by: alex
tags: []
touches: [flai/internal/conflictmark, flai/internal/check, flai/cmd/branch.go, flai/cmd/accept_conflict_test.go, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 2123
  models:
    - model: claude-haiku-4-5-20251001
      input: 178
      output: 5802
      cache_read: 1106506
      cache_write: 70115
      cost: 0.2275
    - model: claude-opus-5-5
      input: 278
      output: 73129
      cache_read: 9255064
      cache_write: 324577
      cost: 5.4998
    - model: claude-sonnet-5
      input: 114
      output: 41029
      cache_read: 2551678
      cache_write: 258062
      cost: 1.566
---
# S-0253 An acceptance merge committed conflict markers to a design document on main, and nothing caught it

## Goal

This story remediates [I-0066](../../../design/issues/I-0066-an-acceptance-merge-committed-conflict-markers-to-a-design-document-on-main-and-nothing-caught-it.md), "An acceptance merge committed conflict markers to a design document on main, and nothing caught it". The issue recommends no solution yet: propose one from its instances before building it.

## Acceptance criteria
- [x] The cause I-0066 describes no longer occurs, with a test that reproduces it where one fits
- [x] I-0066 is closed with `flai issue close I-0066 --reason` saying what fixed it

## Tasks
- T-0858 flai check reports a conflict marker left in markdown as an error
- T-0859 flai accept refuses a story branch that adds a conflict marker

## Notes
