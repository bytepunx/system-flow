---
id: I-0402
title: flai serve stops the dev server of a story in a state a project's workflow adds, because StoryLeft reads the baseline workflow's kinds
class: defect
status: open
count: 1
cost: 10m
first_reported: 2026-11-17T11:35:21Z
last_reported: 2026-11-17T11:35:21Z
updated: 2026-11-17T11:35:21Z
---

# I-0402 flai serve stops the dev server of a story in a state a project's workflow adds, because StoryLeft reads the baseline workflow's kinds

## Description
flai serve stops the dev server of a story in a state a project's workflow adds, because StoryLeft reads the baseline workflow's kinds

## Instances

### 2026-11-17T11:35:21Z
Story: S-0756.
S-0756's end-to-end test (flai/cmd/workflow_e2e_test.go) adds testing, of the active kind, between in-progress and review in a project rendered from the template. serve.StoryLeft (flai/internal/serve/components.go, wf := baselineFlow) answers true for a story in testing, since testing is no baseline state and isUnderway is false, so its caller stops the story's dev server while the story is still active work. Expected false, as for in-progress. The function's doc says it uses the baseline's kinds because the story is all a look reads; it needs the project's workflow passed in (ComponentsOptions). The story's agent itself is kept: flai serve agent restart starts it in testing.

## Remediation
