---
id: T-2066
type: task
nature: feature
title: The tests that reproduce I-0154 name it, and the story records that S-0450 and S-0453 removed its cause
status: done
parent: S-0470
owner: alex
created: 2026-10-23T06:10:01Z
updated: 2026-10-23T06:10:45Z
transitions:
  - to: ready
    at: 2026-10-23T06:10:05Z
    by: agent-S-0470
  - to: in-progress
    at: 2026-10-23T06:10:06Z
    by: agent-S-0470
  - to: done
    at: 2026-10-23T06:10:45Z
    by: agent-S-0470
stream: S-0470
tags: []
touches: [flai/internal/mcpserver/task_test.go, flai/internal/taskdone/taskdone_test.go, wip/kanban/stories/S-0470-task-done-s-answer-repeats-every-open-thread-not-on-the-agent-s-story-so-it-overflows-a-tool-result.md]
usage:
  source: log
  seconds: 39
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 12
      output: 2027
      cache_read: 623847
      cache_write: 28405
      cost: 0.3926
---

# T-2066 The tests that reproduce I-0154 name it, and the story records that S-0450 and S-0453 removed its cause

## Work

I-0154's cause is gone already. S-0450 (ADR-0179, flai v1.67.1) made `task_done` and `story_start` read the inbox story-scoped and fit each answer to 40,000 bytes. S-0453 (ADR-0181, flai v1.67.6) made a thread on another item read `awaiting: operator`, so it is neither listed nor counted. The host's flai is 1.67.13.

Name I-0154 in the comments of the tests that reproduce its instance (the operator's threads on other stories, epics, and a design folder, none on the story): `TestTaskDoneCarriesNoneOfTheOperatorsThreads` in `flai/internal/mcpserver/task_test.go` and `TestRunCarriesNoneOfTheOperatorsThreads` in `flai/internal/taskdone/taskdone_test.go`. Record the finding in the story's notes. No code changes: none is needed. Waits for no task; the story's only one.

## Done when

- Both tests name I-0154 and pass under `flai test`.
- The story's notes say what removed the cause and where it is tested.

## Notes
