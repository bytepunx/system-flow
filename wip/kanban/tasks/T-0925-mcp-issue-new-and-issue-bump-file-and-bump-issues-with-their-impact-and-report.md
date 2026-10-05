---
id: T-0925
type: task
nature: feature
title: MCP issue_new and issue_bump file and bump issues with their impact and report
status: backlog
parent: S-0224
owner: alex
created: 2026-10-05T05:45:02Z
updated: 2026-10-05T05:45:02Z
transitions: []
stream: S-0224
tags: [flai]
touches: [flai/internal/mcpserver/issues.go, flai/internal/mcpserver/issues_test.go, flai/internal/mcpserver/folder.go]
after: [T-0918]
---
# T-0925 MCP issue_new and issue_bump file and bump issues with their impact and report

## Work

Waits for T-0918: the tools call the `issues` functions that task gives the impact and the report link.

- Add the MCP tools `issue_new` (title, class, cost, note, story, the three impact inputs, evidence, report) and `issue_bump` (id, cost, note, story, the impact inputs, evidence, report), registered beside `issue_story`, so an analyzer on a harness without a shell files its findings the same way. Each returns the issue's ID, path, and whether it was new or bumped, refuses what the CLI refuses, regenerates `design/issues/summary.md`, and commits nothing, as `issue_story` does.
- Describe both in the tool list so that the descriptions say the analyzer's use: an actionable finding with its impact, deduplicated against open issues.

## Done when

- `issue_new` with a report bumps an open issue of the same title rather than duplicating it; `issue_bump` links a report and updates the impact
- Tests in `flai/internal/mcpserver/issues_test.go` cover a new issue, a bumped duplicate, and a refusal; `scripts/flai-test.sh` passes

## Notes
