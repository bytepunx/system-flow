---
id: T-3182
type: task
nature: feature
title: flaiover's API routes read the agent images, the jobs, and a job's log, and post a retry
status: done
parent: S-0757
owner: alex
created: 2026-11-12T16:03:46Z
updated: 2026-11-13T16:47:35Z
transitions:
  - to: ready
    at: 2026-11-13T16:40:48Z
    by: agent-S-0757
  - to: in-progress
    at: 2026-11-13T16:40:48Z
    by: agent-S-0757
  - to: done
    at: 2026-11-13T16:47:35Z
    by: agent-S-0757
stream: S-0757
tags: [dashboard]
touches: [flaiover/src/routes/api/images/+server.ts, flaiover/src/routes/api/images/retry/+server.ts, flaiover/src/lib/images.ts, flaiover/src/lib/server/images.ts, flaiover/src/routes/api/images/images.test.ts]
after: [T-3181]
usage:
  source: log
  seconds: 407
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 31
      output: 12025
      cache_read: 1801598
      cache_write: 67738
      cost: 0.9974
---
# T-3182 flaiover's API routes read the agent images, the jobs, and a job's log, and post a retry

## Work

Add `flaiover/src/routes/api/images/+server.ts`. Its GET calls `images.get`, or `images.log` with `job` and `offset`. Add `retry/+server.ts`, whose POST calls `images.ensure`. Both follow `flaiover/src/routes/api/kluster/log/+server.ts` and `.../api/kluster/deploy/+server.ts`, passing a refusal's words through unchanged (ADR-0016).

It waits for T-3181, whose methods it calls.

## Done when

- prettier, eslint, and svelte-check pass on both files.
- A vitest over a fake host channel covers a read, a log from an offset, a retry, and a refused retry.

## Notes
