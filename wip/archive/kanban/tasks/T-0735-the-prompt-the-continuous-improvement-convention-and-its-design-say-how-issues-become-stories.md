---
id: T-0735
type: task
nature: feature
title: The prompt, the continuous-improvement convention, and its design say how issues become stories
status: done
parent: S-0198
owner: arobson
created: 2026-10-03T01:28:26Z
updated: 2026-10-03T02:11:03Z
transitions:
  - to: ready
    at: 2026-10-03T01:28:52Z
    by: agent-S-0198
  - to: in-progress
    at: 2026-10-03T02:09:32Z
    by: agent-S-0198
  - to: done
    at: 2026-10-03T02:11:03Z
    by: agent-S-0198
stream: S-0198
tags: []
touches: [flai/internal/harness, design/conventions/continuous-improvement.md, template/root/design/conventions/continuous-improvement.md, design/system/continuous-improvement.md, design/system/conventions.md, docs/users/conventions.md]
after: [T-0734]
usage:
  source: log
  seconds: 91
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 49
      output: 16605
      cache_read: 2837785
      cache_write: 58164
      cost: 1.2462
---
# T-0735 The prompt, the continuous-improvement convention, and its design say how issues become stories

## Work

- `harness.Prompt` tells the story's agent that `flai issue new` and `bump` record its story, and that before review it makes a story with `flai issue story` (or the MCP tool `issue_story`) for each issue it recorded or bumped that no open story links.
- The continuous-improvement baseline in `template/root/design/conventions/` and its copy in `design/conventions/` say the same, and that the operator can make the rest from the story's review page.
- `design/system/continuous-improvement.md` describes the story recorded per instance, the review page's option, `flai issue story`, and the check.

Waits for the fifth task: it describes what every earlier task built.

## Done when

- The harness test covers the new prompt lines.
- `go test -race -short ./internal/harness/` passes.

## Notes
