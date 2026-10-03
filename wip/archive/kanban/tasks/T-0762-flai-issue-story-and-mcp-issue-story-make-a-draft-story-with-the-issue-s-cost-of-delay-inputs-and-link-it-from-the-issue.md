---
id: T-0762
type: task
nature: improvement
title: flai issue story and MCP issue_story make a draft story with the issue's cost of delay inputs and link it from the issue
status: done
parent: S-0203
owner: alex
created: 2026-10-03T18:07:45Z
updated: 2026-10-03T18:21:22Z
transitions:
  - to: ready
    at: 2026-10-03T18:08:34Z
    by: agent-S-0203
  - to: in-progress
    at: 2026-10-03T18:16:35Z
    by: agent-S-0203
  - to: done
    at: 2026-10-03T18:21:22Z
    by: agent-S-0203
stream: S-0203
tags: []
touches: [flai/cmd/issue.go, flai/cmd/issue_test.go, flai/internal/mcpserver/issues.go, flai/internal/mcpserver/issues_test.go]
after: [T-0760, T-0761]
usage:
  source: log
  seconds: 287
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 50
      output: 18995
      cache_read: 2058110
      cache_write: 74105
      cost: 1.2464
---

# T-0762 flai issue story and MCP issue_story make a draft story with the issue's cost of delay inputs and link it from the issue

## Work

`flai issue story` (and through it the dashboard's issue-to-story flow, which runs it) and MCP `issue_story` pass `ForStory`'s draft flag and cost of delay to `itemnew.Create`, and link the story from the issue with `LinkStory` in the hook, committing the issue with the story under `--autocommit` when it was read from the main checkout. Their help text and description say so. Waits for T-0760 and T-0761, whose functions it calls.

## Done when

- `TestIssueStory` and the MCP issue_story tests check `draft: true`, `cost_of_delay.inputs.time_lost_per_cycle` with `by: flai`, the Notes derivation, and the issue's Remediation line; the autocommit test checks the issue is in the commit.
- `go test -race -short ./cmd/ -run Issue` and `./internal/mcpserver/ -run Issue` pass.

## Notes
