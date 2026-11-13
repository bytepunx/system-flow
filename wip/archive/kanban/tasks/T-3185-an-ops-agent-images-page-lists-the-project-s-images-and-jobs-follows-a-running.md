---
id: T-3185
type: task
nature: feature
title: "An Ops › Agent images page lists the project's images and jobs, follows a running job, and retries a failed one"
status: done
parent: S-0757
owner: alex
created: 2026-11-12T16:04:07Z
updated: 2026-11-13T16:59:47Z
transitions:
  - to: ready
    at: 2026-11-13T16:52:54Z
    by: agent-S-0757
  - to: in-progress
    at: 2026-11-13T16:52:54Z
    by: agent-S-0757
  - to: done
    at: 2026-11-13T16:59:47Z
    by: agent-S-0757
stream: S-0757
tags: [dashboard]
touches: [flaiover/src/routes/ops/images/+page.svelte, flaiover/src/lib/sitemenu.ts, flaiover/src/lib/sitemenu.test.ts, flaiover/src/routes/ops/images/images.svelte.test.ts]
after: [T-3182, T-3183]
usage:
  source: log
  seconds: 413
  estimated: true
  priced_by: harness
  models:
    - model: claude-opus-5-5
      input: 56
      output: 21454
      cache_read: 3214306
      cache_write: 120854
      cost: 1.7795
---
# T-3185 An Ops › Agent images page lists the project's images and jobs, follows a running job, and retries a failed one

## Work

Add `flaiover/src/routes/ops/images/+page.svelte`:

- a table of the project's agent images, their source, and their state, from T-3182's GET;
- an `ImageJob` (T-3183) for the running job and for each last ended one;
- each job's output read from its offset as it grows, as the Deployments page reads a step's log;
- Retry posting to T-3182's retry route, with a refusal shown in its words;
- a job named in the page's address shown first, so that T-3184's link from a story's agent lands on it.

Add the page to the Ops group in `flaiover/src/lib/sitemenu.ts`, beside Klusters and Deployments, and to `sitemenu.test.ts`.

It waits for T-3182 and T-3183, whose route and component it uses.

## Done when

- vitest passes on the page and the menu.
- prettier, eslint, and svelte-check pass.

## Notes
