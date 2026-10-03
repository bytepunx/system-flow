---
id: T-0731
type: task
nature: feature
title: flai issue records the story it runs for, lists issues per story, and makes a story from an issue
status: done
parent: S-0198
owner: arobson
created: 2026-10-03T01:28:18Z
updated: 2026-10-03T01:46:28Z
transitions:
  - to: ready
    at: 2026-10-03T01:28:52Z
    by: agent-S-0198
  - to: in-progress
    at: 2026-10-03T01:34:55Z
    by: agent-S-0198
  - to: done
    at: 2026-10-03T01:46:28Z
    by: agent-S-0198
stream: S-0198
tags: []
touches: [flai/cmd/issue.go, flai/cmd/issue_test.go, docs/users/flai.md, docs/users/flai-reference.md, design/system/flai-cli.md]
after: [T-0730]
usage:
  source: log
  seconds: 693
  estimated: true
  models:
    - model: claude-opus-5-5
      input: 120
      output: 758
      cache_read: 5463689
      cache_write: 159937
      cost: 2.3331
---
# T-0731 flai issue records the story it runs for, lists issues per story, and makes a story from an issue

## Work

In `flai/cmd/issue.go`:

- `flai issue new` and `bump` take `--story`. Without it, the story is `FLAI_STORY`, else the `S-nnnn` in `FLAI_AGENT` (`agent-S-0198`), else the story whose worktree the command runs in (`story/S-nnnn` branch). Outside a story nothing is recorded.
- `flai issue list --story S-nnnn` lists the issues whose instances name that story. `--json` gives each issue its `stories` and the open story that links it, if any.
- `flai issue story I-nnnn [--epic E-nnnn] [--json]` makes a backlog story from `issues.ForStory`, with the project's default agent, as `flai story new` makes one. It records the story under the issue's `## Remediation` and regenerates `summary.md`. An issue that is closed, or already has an open story linking it, is refused, naming that story.

Docs in the same change: the issue section of `docs/users/flai.md`, the `flai issue` row of `design/system/flai-cli.md`, and `docs/users/flai-reference.md` regenerated with `make flai-reference`.

Waits for the first task, whose `issues` functions it calls.

## Done when

- Tests in `flai/cmd` cover the story resolution order, `list --story`, and `issue story` for a defect and an efficiency issue, plus its refusals.
- `go test -race -short ./cmd/ -run Issue` passes.

## Notes
