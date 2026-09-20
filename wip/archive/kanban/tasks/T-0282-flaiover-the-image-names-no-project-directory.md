---
id: T-0282
type: task
nature: feature
title: "flaiover: the image names no project directory"
status: done
parent: S-0077
owner: alex
created: 2026-09-20T12:39:12Z
updated: 2026-09-20T12:44:47Z
transitions:
  - to: ready
    at: 2026-09-20T12:44:47Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T12:44:47Z
    by: system-flow
  - to: done
    at: 2026-09-20T12:44:47Z
    by: system-flow
stream: S-0077
tags: []
---
# T-0282 flaiover: the image names no project directory

## Work
PROJECT_DIR leaves the Dockerfile and the startup log; repo.ts says what the root still means (tests and development only). Anything else in the image that existed for the mount goes.

## Done when
- make flaiover-test and make flaiover-build pass

## Notes
