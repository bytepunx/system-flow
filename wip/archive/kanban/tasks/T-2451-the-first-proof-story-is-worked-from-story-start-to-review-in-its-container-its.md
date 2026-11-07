---
id: T-2451
type: task
nature: experiment
title: The first proof story is worked from story_start to review in its container, its usage measured and its card's status shown as for a host run, and each limitation met is an issue
status: done
parent: S-0565
owner: alex
created: 2026-11-03T01:24:49Z
updated: 2026-11-07T04:47:11Z
transitions:
  - to: ready
    at: 2026-11-07T04:45:15Z
    by: system-flow
  - to: in-progress
    at: 2026-11-07T04:45:16Z
    by: system-flow
  - to: done
    at: 2026-11-07T04:46:20Z
    by: agent-S-0565
stream: S-0565
tags: [docker]
touches: [design/experiments/S-0565-a-story-worked-end-to-end-by-flai-serve-on-harness-container.md, design/issues/summary.md, design/issues/I-0228-claude-code-drops-any-call-to-flai-s-http-mcp-server-that-answers-after-60-seconds-so-wait-for-events-wait-for-work-verify-and-a-slow-task-done-fail-for-an-agent-in-a-container.md]
after: [T-2449]
usage:
  source: log
  seconds: 64
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 11
      output: 1715
      cache_read: 1882662
      cache_write: 13512
      cost: 0.519
---
# T-2451 The first proof story is worked from story_start to review in its container, its usage measured and its card's status shown as for a host run, and each limitation met is an issue

## Work

Waits for T-2450: the stop and the restart are tried early in the same run, and their limitations are filed here with the run's.

Hold `wait_for_events` until the first proof story reaches review, or its agent fails past its automatic restarts. Then record in the results document, under its results:

- the run: each `agent` event with its `why`, the container's name, the image it ran, and how long the start took against a host start;
- `story_start` to review: the tasks closed with `task_done`, the close-out job's verdict (`flai verify S-nnnn --last`), and every step that needed the host;
- the usage: `flai show S-nnnn` after the move to review, its `seconds`, its models, and `priced_by`, against what a host run of a story of its size records;
- the card: the status the dashboard's card and the story's page show while the agent works, waits, and ends, against a host run's.

File each limitation met here and each T-2450 named with `flai issue new`, naming in it the story of E-0029 that carries it on, or none; leave I-0253, I-0254, and I-0259 to their stories and bump them with `flai issue bump` when met again.

## Done when

- The first proof story is in review, worked by its agent in the container from `story_start` on.
- The results document records the run, the usage, and the card's status as above.
- Each limitation met in the run is an issue under `design/issues`, named by ID in the results document.

## Notes
