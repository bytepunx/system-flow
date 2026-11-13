---
id: T-3201
type: task
nature: feature
title: workitem takes its states and transition map from the baseline workflow, and a test hook opens a repository with another workflow
status: backlog
parent: S-0737
owner: alex
created: 2026-11-12T16:06:10Z
updated: 2026-11-12T16:06:10Z
transitions: []
stream: S-0737
tags: [cli]
touches: [flai/internal/workitem/item.go, flai/internal/workitem/rules.go, flai/internal/workitem/transition.go, flai/internal/workitem/repo.go, flai/internal/workitem/rules_test.go]
after: [T-3195]
---
# T-3201 workitem takes its states and transition map from the baseline workflow, and a test hook opens a repository with another workflow

## Work

Waits for T-3195: it reads the baseline that task embeds. Give `Repo` the workflow it runs, the baseline by default, and a way for tests to open a repository with another workflow. Make `States`, the check that a status is a state, `allowed`, `Follows`, and `allowedFor` read it. Keep the six constants as the baseline's names. Add a test that opens a repository with a workflow where `review` is renamed and a state sits between `in-progress` and `review`, and moves a story through it.

## Done when

Every existing test in `flai/...` passes unchanged, and the new test moves a story through the other workflow.

## Notes
