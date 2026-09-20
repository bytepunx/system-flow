---
id: T-0287
type: task
nature: feature
title: push.run pushes pending acceptances and publishes the template, and every host action is journalled
status: done
parent: S-0078
owner: alex
created: 2026-09-20T13:06:48Z
updated: 2026-09-20T13:15:05Z
transitions:
  - to: ready
    at: 2026-09-20T13:15:04Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T13:15:05Z
    by: system-flow
  - to: done
    at: 2026-09-20T13:15:05Z
    by: system-flow
stream: S-0078
tags: []
---
# T-0287 push.run pushes pending acceptances and publishes the template, and every host action is journalled

## Work
flai push --pending gains --publish: after the push, each template component whose publish remote does not match is published as flai template push does, never forced. push.run is that command as a host action. A refusal (the remote moved, the push failed) is a typed error with the reason. Every host action, accept with push included, is written to a journal beside the serve state: when, what, which project, asked by whom, the request, the outcome. flai serve journal reads it.

## Done when
- Tests with a real scratch remote: pushed with tags, refused when the remote moved, nothing pending, disabled; journal entries for each
- make flai-test passes

## Notes
