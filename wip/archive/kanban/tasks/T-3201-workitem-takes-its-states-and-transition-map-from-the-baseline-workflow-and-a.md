---
id: T-3201
type: task
nature: feature
title: workitem takes its states and transition map from the baseline workflow, and a test hook opens a repository with another workflow
status: done
parent: S-0737
owner: alex
created: 2026-11-12T16:06:10Z
updated: 2026-11-14T23:32:46Z
transitions:
  - to: ready
    at: 2026-11-14T23:24:03Z
    by: system-flow
  - to: in-progress
    at: 2026-11-14T23:24:03Z
    by: system-flow
  - to: done
    at: 2026-11-14T23:32:46Z
    by: agent-S-0737
stream: S-0737
tags: [cli]
touches: [flai/internal/workitem/item.go, flai/internal/workitem/rules.go, flai/internal/workitem/repo.go, flai/internal/workitem/rules_test.go]
after: [T-3195]
usage:
  source: log
  seconds: 523
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 27
      output: 13307
      cache_read: 3239682
      cache_write: 74545
      cost: 1.4625
---
# T-3201 workitem takes its states and transition map from the baseline workflow, and a test hook opens a repository with another workflow

## Work

Waits for T-3195: it reads the baseline that task embeds. Give `Repo` the workflow it runs, the baseline by default, and a way for tests to open a repository with another workflow. Make `States`, the check that a status is a state, `allowed`, `Follows`, and `allowedFor` read it. Keep the six constants as the baseline's names. Add a test that opens a repository with a workflow where `review` is renamed and a state sits between `in-progress` and `review`, and moves a story through it.

## Done when

Every existing test in `flai/...` passes unchanged, and the new test moves a story through the other workflow.

## Notes
