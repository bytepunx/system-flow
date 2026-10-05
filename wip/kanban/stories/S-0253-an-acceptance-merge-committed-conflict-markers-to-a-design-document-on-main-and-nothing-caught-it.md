---
id: S-0253
type: story
nature: remediation
title: An acceptance merge committed conflict markers to a design document on main, and nothing caught it
status: in-progress
owner: alex
created: 2026-10-03T20:06:56Z
updated: 2026-10-05T03:20:57Z
transitions:
  - to: ready
    at: 2026-10-04T21:42:57Z
    by: alex
  - to: in-progress
    at: 2026-10-05T03:10:48Z
    by: agent-S-0253
tags: []
touches: [flai/internal/conflictmark, flai/internal/check, flai/cmd/branch.go, flai/cmd/accept_conflict_test.go, design/system/flai-cli.md, design/system/workflow.md, docs/users/flai.md]
agent:
  harness: claude-code
  model: claude-opus-5-5
  config:
    effort: high
usage:
  source: log
  seconds: 1254
  models:
    - model: claude-haiku-4-5-20251001
      input: 178
      output: 5802
      cache_read: 1106506
      cache_write: 70115
      cost: 0.2275
    - model: claude-opus-5-5
      input: 214
      output: 59273
      cache_read: 7443039
      cache_write: 259329
      cost: 4.338
    - model: claude-sonnet-5
      input: 38
      output: 11540
      cache_read: 776950
      cache_write: 67269
      cost: 0.439
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
