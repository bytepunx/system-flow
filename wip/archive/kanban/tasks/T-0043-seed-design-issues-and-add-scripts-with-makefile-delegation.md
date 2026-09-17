---
id: T-0043
type: task
nature: feature
title: Seed design/issues and add scripts/ with Makefile delegation
status: done
parent: S-0026
owner: alex
created: 2026-09-15T22:37:53Z
updated: 2026-09-15T22:42:25Z
transitions:
  - to: ready
    at: 2026-09-15T22:42:25Z
    by: agent
  - to: in-progress
    at: 2026-09-15T22:42:25Z
    by: agent
  - to: done
    at: 2026-09-15T22:42:25Z
    by: agent
stream: S-0026
tags: [conventions, tooling]
---

# T-0043 Seed design/issues and add scripts/ with Makefile delegation

## Work
Seed design/issues with the friction observed so far (six issues with counts and commit-derived timestamps) and summary.md; add scripts/ (env, flai wrapper, build, test, snapshot, check, template-test, install-tools) and make the Makefile delegate to them; template gets scripts/README, check.sh, Makefile delegation, and an issues skeleton.

## Done when
Scripts run end to end; summary lists every open issue.

## Notes
