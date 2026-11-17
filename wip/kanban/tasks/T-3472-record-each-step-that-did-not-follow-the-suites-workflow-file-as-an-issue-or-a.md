---
id: T-3472
type: task
nature: feature
title: Record each step that did not follow the suites workflow file as an issue or a story under E-0034, and say in workflow.md that ADR-0312's outcome holds
status: backlog
parent: S-0815
owner: alex
created: 2026-11-16T23:50:37Z
updated: 2026-11-16T23:50:37Z
transitions: []
stream: S-0815
tags: [flai]
touches: [design/system/workflow.md]
after: [T-3463, T-3466, T-3468]
---
# T-3472 Record each step that did not follow the suites workflow file as an issue or a story under E-0034, and say in workflow.md that ADR-0312's outcome holds

## Work

It waits for T-3463, T-3466, and T-3468 because it records what their tests found.

For each step that did not follow the file in those tests, such as a close-out that ran a slow suite, a job that ran one twice, a cluster step that held a story, or a dashboard read that lacked the state, record an issue with `flai issue new`, or draft a story under E-0034 when the fix is a deliverable of its own, naming the test that shows it. In `design/system/workflow.md`, under the workflow file's section, say that a project can run only its fast suites before review and its slow suites in an integration state, cite the test, and link each issue or story recorded.

## Done when

Every failure the three tests found is an issue or a story under E-0034, and `design/system/workflow.md` says what the proof showed.

## Notes
