---
id: T-3259
type: task
nature: feature
title: flai check runs the workflow findings on the project's own file and advises a project with none, if the ADR asks
status: backlog
parent: S-0755
owner: alex
created: 2026-11-12T16:10:11Z
updated: 2026-11-12T16:10:11Z
transitions: []
stream: S-0755
tags: [cli]
touches: [flai/internal/check/check_workflow.go, flai/internal/check/workflow_test.go]
after: [T-3253]
---
# T-3259 flai check runs the workflow findings on the project's own file and advises a project with none, if the ADR asks

## Work

Waits for T-3253: the check reads the file that task loads. Make the `workflow.` findings run on the project's own file when it has one, naming the file and the line, and add the advice for a project with no file, as `roles.baseline` advises one with no role file, if the ADR asks for it. Test a project file that fails a finding and one that passes.

## Done when

`workflow_test.go` gives the findings on a project's own file and the advice where the ADR asks.

## Notes
