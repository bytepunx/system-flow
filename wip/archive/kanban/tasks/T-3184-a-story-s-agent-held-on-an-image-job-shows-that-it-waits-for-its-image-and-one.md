---
id: T-3184
type: task
nature: feature
title: A story's agent held on an image job shows that it waits for its image, and one failed by a job shows the cause, each linked to the job
status: done
parent: S-0757
owner: alex
created: 2026-11-12T16:04:02Z
updated: 2026-11-13T16:52:10Z
transitions:
  - to: ready
    at: 2026-11-13T16:40:49Z
    by: agent-S-0757
  - to: in-progress
    at: 2026-11-13T16:40:49Z
    by: agent-S-0757
  - to: done
    at: 2026-11-13T16:52:10Z
    by: agent-S-0757
stream: S-0757
tags: [dashboard]
touches: [flaiover/src/lib/components/StoryAgent.svelte, flaiover/src/lib/components/StoryAgent.svelte.test.ts, flaiover/src/lib/activity.ts, flaiover/src/routes/ops/images/+page.svelte]
after: [T-3181]
usage:
  source: log
  seconds: 681
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 61
      output: 23213
      cache_read: 3477867
      cache_write: 130763
      cost: 1.9254
---
# T-3184 A story's agent held on an image job shows that it waits for its image, and one failed by a job shows the cause, each linked to the job

## Work

In `flaiover/src/lib/components/StoryAgent.svelte`, read the run's image-job fields that S-0751 records in `agents.json` and the host API passes on:

- a run waiting for its image shows "waiting for its image" and the job's step, with a link to the job on the Ops page;
- a run failed by a job shows the job's cause, with the same link, in place of the bare refusal text that `act` shows today (lines 115 to 128).

Test both in `StoryAgent.svelte.test.ts`.

It waits for T-3181, which carries the fields to the dashboard. It shares no file with T-3182 or T-3183, so the three run together.

## Done when

- vitest passes on `StoryAgent.svelte`, with the two new cases and the old ones unchanged.
- prettier, eslint, and svelte-check pass.

## Notes
