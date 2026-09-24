---
id: T-0398
type: task
nature: feature
title: Verify the by-hand walkthrough
status: done
parent: S-0016
owner: agent
created: 2026-09-24T08:01:48Z
updated: 2026-09-24T08:04:36Z
transitions:
  - to: ready
    at: 2026-09-24T08:03:31Z
    by: agent-S-0016
  - to: in-progress
    at: 2026-09-24T08:03:32Z
    by: agent-S-0016
  - to: done
    at: 2026-09-24T08:04:36Z
    by: agent-S-0016
stream: S-0016
tags: []
touches: [docs/users/conventions.md]
---
# T-0398 Verify the by-hand walkthrough

## Work
Follow the walkthrough in docs/users/conventions.md literally in a scratch directory, without flai new or flai import, then run flai check --strict on the result, move a story through its states with flai, and fix every gap the guide left.

## Done when
A repository built only from the guide passes flai check --strict and flai accepts work on it.

## Notes
