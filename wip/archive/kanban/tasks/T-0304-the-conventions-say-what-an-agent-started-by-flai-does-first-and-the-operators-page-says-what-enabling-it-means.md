---
id: T-0304
type: task
nature: feature
title: The conventions say what an agent started by flai does first, and the operators' page says what enabling it means
status: done
parent: S-0079
owner: alex
created: 2026-09-20T15:47:56Z
updated: 2026-09-20T16:01:12Z
transitions:
  - to: ready
    at: 2026-09-20T15:59:45Z
    by: system-flow
  - to: in-progress
    at: 2026-09-20T15:59:46Z
    by: system-flow
  - to: done
    at: 2026-09-20T16:01:12Z
    by: system-flow
stream: S-0079
tags: []
---
# T-0304 The conventions say what an agent started by flai does first, and the operators' page says what enabling it means

## Work
session-start baseline, in template/ first and identical here: FLAI_STORY names the story flai started the session for. Operators: what enabling means (a holder of the dashboard token, or anyone who can move a story to ready, starts your command on your machine), how to set the command, the attending rule, the journal. Design pages.

## Done when
- The markdown lint and flai check --strict pass after the last edit, in the worktree and the main checkout; the template renders and checks

## Notes
