---
id: T-2787
type: task
nature: remediation
title: story_start's refusal of a story held no-touches by its own empty claim says that declaring its touches clears the hold
status: done
parent: S-0633
owner: alex
created: 2026-11-09T22:01:36Z
updated: 2026-11-09T23:14:57Z
transitions:
  - to: ready
    at: 2026-11-09T22:21:11Z
    by: system-flow
  - to: in-progress
    at: 2026-11-09T22:21:12Z
    by: system-flow
  - to: done
    at: 2026-11-09T23:14:57Z
    by: agent-S-0633
stream: S-0633
tags: [flai, storystart]
touches: [flai/internal/storystart/start.go, flai/internal/storystart/start_test.go]
usage:
  source: log
  seconds: 3225
  estimated: true
  priced_by: estimate
  models:
    - model: claude-opus-5-5
      input: 35
      output: 283
      cache_read: 5486279
      cache_write: 52099
      cost: 2.0251
---
# T-2787 story_start's refusal of a story held no-touches by its own empty claim says that declaring its touches clears the hold

## Work

`refuse` in `flai/internal/storystart/start.go` answers every hold with the hold's reason and "nothing was changed, so pull the first ready story the inbox lists that is not held". That advice is wrong for the agent flai serve started for this one story when the story is held `no-touches` because it declares no touches itself. That agent cannot pull another story, and what clears the hold is its own `flai touches`. In I-0301, agent-S-0568 found that out by reading the guard and the issue.

When the hold's code is `workitem.HoldNoTouches` and the story's own claim is empty, end the refusal with what to do instead. Tell the agent to declare the paths the story changes with `flai touches <S-nnnn> <path>...`, or `item_edit` `touches`, and to call `story_start` again. Keep the current advice for every other hold, and for a no-touches hold that only another story's empty claim causes. Add a test in `flai/internal/storystart/start_test.go` next to `TestStartRefusesAStoryHeldByAnOverlap`, which has no no-touches case yet. It covers both no-touches cases. If `flai/internal/mcpserver/story_start_test.go` or `flai/cmd/story_start_test.go` assert the old text, update them; `task_done` widens the touches.

Waits for no task: it shares no path with T-2786.

## Done when

- `story_start` on a ready story with no touches and no tasks, while another story is in progress, refuses it, naming `flai touches` and `story_start` again
- A story held no-touches only by another story's empty claim keeps the current advice
- The new tests in `flai/internal/storystart/start_test.go` pass
- `flai test flai/internal/storystart flai/internal/mcpserver flai/cmd` passes

## Notes
