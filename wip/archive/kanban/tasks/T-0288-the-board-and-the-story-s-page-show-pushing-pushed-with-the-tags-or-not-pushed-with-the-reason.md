---
id: T-0288
type: task
nature: feature
title: The board and the story's page show pushing, pushed with the tags, or not pushed with the reason
status: done
parent: S-0078
owner: alex
created: 2026-09-20T13:06:48Z
updated: 2026-09-20T13:19:07Z
transitions:
  - to: ready
    at: 2026-09-20T13:15:05Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T13:15:05Z
    by: system-flow
  - to: done
    at: 2026-09-20T13:19:07Z
    by: system-flow
stream: S-0078
tags: []
---
# T-0288 The board and the story's page show pushing, pushed with the tags, or not pushed with the reason

## Work
project.info says which host actions are enabled. Acceptance shows the push as its last step when enabled; the standing notice offers a push when enabled and keeps the command to run by hand; a refusal shows its reason. Nothing in the dashboard can enable an action.

## Done when
- Component and route tests; make flaiover-test and make flaiover-build pass

## Notes
